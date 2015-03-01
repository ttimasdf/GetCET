#!/usr/bin/env python
# coding=utf-8

import requests
import CetConfig
from ctypes import CDLL, c_char, c_int, c_long, byref, pointer, \
    create_string_buffer, Structure, Union, sizeof

DES_cblock = c_char * 8
DES_LONG = c_int


class ks(Union):
    _fields_ = [
        ('cblock', DES_cblock),
        ('deslong', DES_LONG * 2)
    ]


class DES_key_schedule(Structure):
    _fields_ = [
        ('ks', ks * 16),
    ]


class CetCrypter(object):

    ticket_number_key = '(YesuNRY'
    request_data_key = '?!btwNP^'

    ENCRYPT = c_int(1)
    DECRYPT = c_int(0)

    def __init__(self, libcrypto_path=None):
        if not libcrypto_path:
            from ctypes.util import find_library
            libcrypto_path = find_library('crypto')
            if not libcrypto_path:
                raise Exception('libcrypto(OpenSSL) not found')

        self.libcrypto = CDLL(libcrypto_path)

        if hasattr(self.libcrypto, 'OpenSSL_add_all_ciphers'):
            self.libcrypto.OpenSSL_add_all_ciphers()

    def process_data(self, indata, key, is_enc=1):
        length = len(indata)

        indata = create_string_buffer(indata, length)
        outdata = create_string_buffer(length)
        n = c_int(0)

        key = DES_cblock(*tuple(key))
        key_schedule = DES_key_schedule()

        self.libcrypto.DES_set_odd_parity(key)
        self.libcrypto.DES_set_key_checked(byref(key), byref(key_schedule))

        self.libcrypto.DES_cfb64_encrypt(byref(indata),
                                         byref(outdata),
                                         c_int(length),
                                         byref(key_schedule), byref(key), byref(n), c_int(is_enc))

        return outdata.raw

    def decrypt_ticket_number(self, ciphertext):
        ciphertext = ciphertext[2:]
        return self.process_data(ciphertext, self.ticket_number_key, 0)

    def encrypt_ticket_number(self, ticket_number):
        ciphertext = self.process_data(
            ticket_number, self.ticket_number_key, 1)
        ciphertext = '\x35\x2c' + ciphertext
        return ciphertext

    def decrypt_request_data(self, ciphertext):
        return self.process_data(ciphertext, self.request_data_key, 0)

    def encrypt_request_data(self, plaintext):
        return self.process_data(plaintext, self.request_data_key, 1)


class CetTicket(object):

    """
        usage:
        ct = CetTicket()
        print ct.find_ticket_number(b'浙江', b'浙江海洋学院', b'XXX', cet_type=2)
    """

    search_url = CetConfig.SEARCH_URL

    def __init__(self):
        self.crypter = CetCrypter()

    def find_ticket_number(self, provice, school, name, examroom='', cet_type=1):
        """
            You can read the `school.json` file to check if your school is supported.
            cet_type:
                    1 ==> cet4
                    2 ==> cet6
        """
        provice_id = CetConfig.PROVICE[provice]
        param_data = b'type=%d&provice=%d&school=%s&name=%s&examroom=%s' % (cet_type,
                                                                            provice_id,
                                                                            school, name, examroom)

        param_data = param_data.decode('utf-8').encode('gb2312')
        encrypted_data = self.crypter.encrypt_request_data(param_data)

        resp = requests.post(url=self.search_url, data=encrypted_data)

        ticket_number = self.crypter.decrypt_ticket_number(resp.content)
        if ticket_number == '':
            raise Exception('Cannot find ticket number.')

        return ticket_number

    def get_score(self, ticket_number, name):
        if isinstance(name, unicode):
            name = name.encode('gb2312')
        else:
            name = name.decode('utf-8').encode('gb2312')

        params_dict = {
            'id': ticket_number,
            'name': name[:4]
        }

        resp = requests.post(url=CetConfig.SCORE_URL,
                             data=params_dict,
                             headers={'Referer': 'http://cet.99sushe.com/'})
        score_data = resp.content.decode('gb2312').encode('utf-8')
        score_data = score_data.split(',')

        score = {
            'name': score_data[6],
            'school': score_data[5],
            'Listening': score_data[1],
            'Reading': score_data[2],
            'Writing': score_data[3],
            'Total': score_data[4]
        }
        return score

if __name__ == '__main__':
    ct = CetTicket()
    print ct.find_ticket_number(b'浙江', b'浙江海洋学院', b'XXX', cet_type=2)
    print ct.get_score('330400XXXXXXXXX', b'XXX')

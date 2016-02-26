FROM python:2.7.10

ADD . /GetCet
WORKDIR /GetCet

RUN pip install -r requirements.txt

ENTRYPOINT ["python", "GetCET.py"]
EXPOSE 8081
CMD ["--address=0.0.0.0", "--port=8081"]

FROM python:2.7.10

ADD . /GetCet
WORKDIR /GetCet

ENV USER_AGENT 高坂穗乃果

RUN pip install -r requirements.txt

ENTRYPOINT ["python", "GetCET.py"]
EXPOSE 8081
CMD ["--address=0.0.0.0", "--port=8081"]

FROM debian:latest

RUN apt-get -y update && apt-get install -y python python-pip

ADD . /GetCet
WORKDIR /GetCet

RUN pip install -r requirements.txt

ENTRYPOINT ["python", "GetCET.py"]
EXPOSE 8081
CMD ["--port=8081"]

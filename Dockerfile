FROM debian:latest

RUN apt-get -y update && apt-get install -y python && pip install -r requirements.txt

ADD . /GetCet
WORKDIR /GetCet

ENTRYPOINT ["python", "./GetCet.py"]
EXPOSE 8081
CMD ["--port=8081"]

FROM selenoid/vnc:chrome_91.0
USER root
RUN apt-get update -qqy && apt-get -qqy install  language-pack-zh-han*  && apt-get -qqy install lame   && apt-get -qqy install ffmpeg
USER selenium
COPY vncproxy /home/selenium/vncproxy
COPY static /home/selenium/static
COPY entrypoint.sh /entrypoint.sh
EXPOSE 8888 5900

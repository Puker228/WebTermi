FROM ubuntu:24.04

RUN apt-get update && apt-get install -y \ 
	locales \
	sudo \
	bash \
	curl \
	wget\
	git \
	nano \
	htop \
	make \
	tar \
	&& localedef -i en_US -c -f UTF-8 -A /usr/share/locale/locale.alias en_US.UTF-8
ENV LANG=en_US.utf8

RUN useradd -m -s /bin/bash dev && \
    usermod -aG sudo dev && \
	echo "dev ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/dev

USER dev
WORKDIR /home/dev

CMD ["bash"]
FROM alpine:latest

WORKDIR /root
COPY xray.sh /root/xray.sh
COPY supervisord.conf /etc/supervisord.conf
RUN set -ex \
	&& apk add --no-cache tzdata ca-certificates supervisor \
	&& mkdir -p /var/log/au /usr/local/share/xray /usr/local/etc/au /usr/local/etc/xray \
	&& chmod +x /root/xray.sh \
	&& /root/xray.sh \
	&& rm -fv /root/xray.sh \
	&& wget -O /usr/local/share/xray/geosite.dat https://github.com/v2fly/domain-list-community/releases/latest/download/dlc.dat \
	&& wget -O /usr/local/share/xray/geoip.dat https://github.com/v2fly/geoip/releases/latest/download/geoip.dat \
	&& wget -O /usr/local/bin/au https://github.com/crossfw/Air-Universe/releases/latest/download/Air-Universe-linux-amd64 \
	&& wget -O /usr/local/etc/xray/xr.json https://raw.githubusercontent.com/crossfw/Air-Universe/master/configs/xray_json/multiIn.json \
	&& chmod +x /usr/local/bin/au

ENV TZ=Asia/Shanghai
CMD ["/usr/bin/supervisord"]

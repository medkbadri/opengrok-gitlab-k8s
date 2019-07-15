FROM opengrok/docker:latest

COPY ./scripts/start-opengrok.sh /scripts/start-opengrok.sh
COPY ./scripts/index-opengrok.sh /scripts/index-opengrok.sh
COPY ./jspf/foot.jspf $CATALINA_TMPDIR/foot.jspf

RUN chmod +x /scripts/start-opengrok.sh /scripts/index-opengrok.sh

CMD ["/scripts/start-opengrok.sh"]
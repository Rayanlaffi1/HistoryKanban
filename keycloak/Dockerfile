FROM quay.io/keycloak/keycloak:26.6

ADD --chown=keycloak:keycloak https://repo1.maven.org/maven2/io/phasetwo/keycloak/keycloak-magic-link/0.75/keycloak-magic-link-0.75.jar /opt/keycloak/providers/keycloak-magic-link.jar
ADD --chown=keycloak:keycloak https://repo1.maven.org/maven2/io/phasetwo/keycloak/keycloak-events/0.62/keycloak-events-0.62.jar /opt/keycloak/providers/keycloak-events.jar

RUN /opt/keycloak/bin/kc.sh build

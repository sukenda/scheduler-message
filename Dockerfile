FROM rabbitmq:3.9-management

ADD plugins/rabbitmq_delayed_message_exchange-3.9.0.ez $RABBITMQ_HOM/plugins/

RUN rabbitmq-plugins enable --offline rabbitmq_delayed_message_exchange
import pika

from .handler import callback
from .settings import broker_settings


def run():
    connection = pika.BlockingConnection(
        pika.ConnectionParameters(
            host=broker_settings.host,
            port=broker_settings.port,
            virtual_host=broker_settings.virtualhost,
            credentials=pika.PlainCredentials(
                username=broker_settings.username,
                password=broker_settings.password.get_secret_value(),
            ),
        )
    )
    channel = connection.channel()
    channel.queue_declare(queue=broker_settings.work_queue, durable=True)
    channel.queue_declare(queue=broker_settings.target_queue, durable=True)

    channel.basic_consume(
        queue=broker_settings.work_queue, auto_ack=True, on_message_callback=callback
    )

    channel.start_consuming()

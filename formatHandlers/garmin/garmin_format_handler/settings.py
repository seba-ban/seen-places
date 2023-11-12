from pydantic import SecretStr
from pydantic_settings import BaseSettings, SettingsConfigDict


class GarminHandlerSettings(BaseSettings):
    local_storage_files_dir: str


class BrokerSettings(BaseSettings):
    host: str
    port: int = 5672
    username: str
    password: SecretStr
    virtualhost: str = "/"
    work_queue: str = "garmin"
    target_queue: str = "points"
    points_buffer_size: int = 10000

    model_config = SettingsConfigDict(env_prefix="BROKER_")


settings = GarminHandlerSettings()
broker_settings = BrokerSettings()

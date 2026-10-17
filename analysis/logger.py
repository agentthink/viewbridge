"""
统一日志管理模块
功能：
1. 根据命令行参数设置日志级别。
2. 支持日志文件按天轮转 (Daily Rotation)。
3. 支持自动清理旧日志 (Retention)。
4. 同时输出到控制台 (Console) 和文件 (File)。
"""

import logging
import logging.handlers
import os
import sys

LOG_DIR = "../configs/logs"
LOG_FILE = "analysis.log"


def setup_logging(level_str: str = "INFO", retention_days: int = 3):
    # 1. 转换级别
    level = getattr(logging, level_str.upper(), logging.INFO)

    # 2. 确保日志目录存在
    if not os.path.exists(LOG_DIR):
        os.makedirs(LOG_DIR)

    log_path = os.path.join(LOG_DIR, LOG_FILE)

    # 3. 创建 Root Logger
    root_logger = logging.getLogger()
    root_logger.setLevel(level)

    # 清除已有的 handlers (避免重复打印)
    root_logger.handlers = []

    # 4. 格式化器
    formatter = logging.Formatter(
        fmt="%(asctime)s | %(levelname)-8s | %(process)d:%(threadName)s | %(message)s",
        datefmt="%Y-%m-%d %H:%M:%S",
    )

    # 5. Handler 1: 控制台输出
    console_handler = logging.StreamHandler(sys.stdout)
    console_handler.setFormatter(formatter)
    root_logger.addHandler(console_handler)

    # 6. Handler 2: 文件输出 (按天轮转)
    # TimedRotatingFileHandler:
    # - when='midnight': 每天午夜切分
    # - interval=1: 每1天
    # - backupCount=retention_days: 保留几天，超出的会被删除
    # - encoding='utf-8': 防止中文乱码
    file_handler = logging.handlers.TimedRotatingFileHandler(
        filename=log_path,
        when="midnight",
        interval=1,
        backupCount=retention_days,
        encoding="utf-8",
    )

    # 设置后缀格式，例如 app.log.2023-12-31
    file_handler.suffix = "%Y-%m-%d.log"
    file_handler.setFormatter(formatter)
    root_logger.addHandler(file_handler)

    # 打印一条初始化日志
    logging.info(
        f"日志系统已初始化: Level={level_str}, Path={log_path}, Retention={retention_days} days"
    )

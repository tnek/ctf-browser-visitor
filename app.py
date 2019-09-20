#!/usr/bin/python
import os
import asyncio
import threading
import json
import logging

from quart import Quart
from quart import request

from xssbot import XSSBotManager

logging.basicConfig(level=logging.INFO)

DEBUG = True
WORKER_POOL_SIZE = 10
HOST = os.environ.get("XSSBOT_HOST", "0.0.0.0")
PORT = os.environ.get("XSSBOT_PORT", 5000)
if DEBUG:
    WORKER_POOL_SIZE = 0
    logging.basicConfig(level=logging.DEBUG)


manager = XSSBotManager(WORKER_POOL_SIZE)


app = Quart(__name__)
loop = asyncio.get_event_loop()
loop.create_task(manager.event_loop())


@app.route("/")
async def index():
    return "asdf"


@app.route("/visit", methods=["GET", "POST"])
async def visit():
    form = await request.args
    job = form.get("job", None)
    if not job:
        return "{}"

    config = json.loads(job)
    return config


if __name__ == "__main__":
    app.run(HOST, PORT)

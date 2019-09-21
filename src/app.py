#!/usr/bin/python
import os
import asyncio
import concurrent.futures
import threading
import json
import logging

from quart import Quart, request

import xssbot
from config import MAX_WORKER_COUNT, REQUIRED_FIELDS

app = Quart(__name__)
logger = logging.getLogger("xssbotapp")

@app.before_serving
async def app_init():
    app.loop = asyncio.get_event_loop()
    await xssbot.xssbot(app.loop)


@app.route("/visit", methods=["GET", "POST"])
async def visit():
    form = request.args
    job = form.get("job", None)
    if not job:
        logger.info("Missing log field")
        return "{}"
    try:
      config = json.loads(job)
      logger.info(config)
      if not all(field in config for field in REQUIRED_FIELDS):
          return '{"status":"fail"}'

      await xssbot.queue_job(config)
      return '{"status":"ok"}'
    except:
      return 400, '{}'


if __name__ == "__main__":
    app.run("0.0.0.0", port=5000)

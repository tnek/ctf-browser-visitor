#!/usr/bin/python
import asyncio
import logging
import sys

from arsenic import get_session, services, browsers

from config import MAX_WORKER_COUNT

workers_sem = asyncio.Semaphore(MAX_WORKER_COUNT)
job_q = None


async def visit(config):
    service = services.Geckodriver()
    browser = browsers.Firefox(**{"moz:firefoxOptions": {"args": ["-headless"]}})
    
    logging.info("Hitting url " + config['url'])
    try:
      async with get_session(service, browser) as session:
          await session.delete_all_cookies()
          await session.get(config['url'])

          for k, c in config.get("cookies", {}).items():
              value = c.get('value', '')
              domain = c.get('domain', None)
              path = c.get('path', "/")
              secure = c.get('secure', False)
              await session.add_cookie(k, value, path=path, domain=domain, secure=secure)

          await session.get(config["url"])
    except Exception as e:
      logging.info("Exception hitting url " + str(config) + " with exception " + e.message)


async def queue_job(job):
    await job_q.put(job)


async def event_loop():
    while True:
        job = await job_q.get()
        async with workers_sem:
            logging.info(job)
            await visit(job)


async def xssbot(loop):
    loop.create_task(event_loop())
    global job_q
    job_q = asyncio.Queue(loop=loop)

#!/usr/bin/python
import asyncio
import logging

from arsenic import get_session, services, browsers

from config import MAX_WORKER_COUNT

workers_sem = asyncio.Semaphore(MAX_WORKER_COUNT)
job_q = asyncio.Queue()


async def visit(config):
    service = services.Geckodriver()
    browser = browsers.Firefox(**{"moz:firefoxOptions": {"args": ["-headless"]}})

    async with get_session(service, browser) as session:
        await session.delete_all_cookies()
        for c in config.get("cookies", {}):
            await session.add_cookie(c, config["cookies"][c])

        await session.get(config["url"])


async def queue_job(job):
    await job_q.put(job)


async def event_loop():
    while True:
        job = await job_q.get()
        async with workers_sem:
            await visit(job)


async def xssbot(loop):
    loop.create_task(event_loop())

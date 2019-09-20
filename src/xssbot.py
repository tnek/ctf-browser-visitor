#!/usr/bin/python
import asyncio
import logging

from arsenic import get_session, services, browsers

job_q = []
job_q_l = asyncio.Lock()


async def visit(config):
    service = services.Geckodriver()
    browser = browsers.Firefox(**{"moz:firefoxOptions": {"args": ["-headless"]}})

    async with get_session(service, browser) as session:
        await session.delete_all_cookies()
        for c in config.get("cookies", {}):
            await session.add_cookie(c, config["cookies"][c])

        await session.get(config["url"])


async def queue_job(job):
    async with job_q_l:
        job_q.append(job)


async def event_loop():
    while True:
        if job_q:
            while job_q:
                async with job_q_l:
                    await visit(job_q.pop())

        await asyncio.sleep(1)


async def xssbot(loop):
    loop.create_task(event_loop())

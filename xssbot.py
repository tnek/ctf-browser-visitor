#!/usr/bin/python
import asyncio
import logging

from selenium import webdriver
from selenium.webdriver.firefox.options import Options


class _XSSBotWorker(object):
    def __init__(self):
        self.options = Options()
        self.options.headless = True
        self.driver = webdriver.Firefox(options=self.options)

    async def visit(self, config):
        self.driver.delete_all_cookies()
        self.driver.implicitly_wait(30)
        self.driver.add_cookie(config["cookies"])

        if self.config["method"] == "GET":
            await self.driver.get(config["url"])


class XSSBotManager(object):
    def __init__(self, pool_size):
        logging.info("Instantiating workers")
        self.workers_l = asyncio.Lock()
        self.workers = [_XSSBotWorker() for i in range(pool_size)]
        self.job_q = []
        self.job_q_l = asyncio.Lock()

    async def queue_job(self, job):
        async with self.job_q_l:
            self.job_q.append(job)

    async def event_loop(self):
        while True:
            if self.job_q:
                async with self.job_q_l:
                    async with self.workers_l:
                        worker = self.workers.pop()
                    await worker.visit(self.job_q.pop())

                async with self.workers_l:
                    self.workers.append(worker)
            await asyncio.sleep(1)

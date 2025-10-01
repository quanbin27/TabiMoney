import asyncio

class AnomalyService:
    def __init__(self, ml_service):
        self.ml_service = ml_service
        self._ready = False

    async def initialize(self):
        await asyncio.sleep(0.01)
        self._ready = True

    async def cleanup(self):
        self._ready = False

    def is_ready(self):
        return self._ready



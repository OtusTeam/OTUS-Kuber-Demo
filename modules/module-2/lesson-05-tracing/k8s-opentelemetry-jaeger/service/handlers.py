import asyncio
import json
import logging
import os
from multiprocessing import Process

import aiofiles
from flask import Flask, render_template
from kubernetes import client, config
from kubernetes.client.rest import ApiException

KUBERNETES_CLUSTER = os.environ.get("KUBERNETES_CLUSTER", "undefined")
NAME_LIST = [
    "cert-manager",
    "trust-manager"
]
STATE_FILE = "state.json"


logger = logging.getLogger("datascraper")
handler = logging.StreamHandler()
handler.setFormatter(
    logging.Formatter("%(asctime)s - %(name)s - %(levelname)s - %(message)s")
)
logger.addHandler(handler)
logger.setLevel(logging.INFO)

app = Flask(__name__)


@app.route("/state")
async def get_info():
    async with aiofiles.open(STATE_FILE, "r") as f:
        data = await f.read()
    state = json.loads(data)
    return state


@app.route("/")
async def index():
    """Render the main page with state data in a Bootstrap template"""
    try:
        async with aiofiles.open(STATE_FILE, "r") as f:
            data = await f.read()
        state_data = json.loads(data)
    except FileNotFoundError:
        state_data = {}

    return render_template("index.html", state_data=state_data)


def flask_process():
    app.run(port=8080, host="0.0.0.0")


def init_state():
    try:
        # Load kubeconfig
        config.load_incluster_config()

        # Create API instance
        api_instance = client.CustomObjectsApi()

        # Get HelmRelease resources
        helm_releases = api_instance.list_cluster_custom_object(
            group="helm.toolkit.fluxcd.io",
            version="v2beta1",
            plural="helmreleases",
        )

        with open(STATE_FILE, "w") as f:
            dict_versions = {}
            for release in helm_releases["items"]:
                name = release["metadata"]["name"]
                if name in NAME_LIST:
                    namespace = release["metadata"]["namespace"]
                    version = (
                        release.get("spec", {})
                        .get("chart", {})
                        .get("spec", {})
                        .get("version", "Unknown")
                    )
                    version = version.replace("v", "")
                    logger.info(
                        f"HelmRelease {name} in namespace {namespace}: version {version}"
                    )
                    dict_versions.update({name: version})
            f.write(json.dumps(dict_versions, indent=4))

    except ApiException as e:
        logger.error(f"Exception when getting HelmRelease: {e}")


if __name__ == "__main__":
    init_state()
    web = Process(target=flask_process, daemon=True)
    web.start()
    web.join()

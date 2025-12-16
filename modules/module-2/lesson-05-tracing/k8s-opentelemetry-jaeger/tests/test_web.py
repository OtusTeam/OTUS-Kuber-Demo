import json
import os

import pytest

import service.handlers


@pytest.fixture
def client():
    with service.handlers.app.test_client() as client:
        yield client


def test_root_page(client):
    test_data = {
        "cert-manager": "1.11.0",
        "ingress-nginx": "4.5.2",
        "prometheus-stack": "50.0.0",
    }
    with open("state.json", "w") as f:
        json.dump(test_data, f)
    try:
        response = client.get("/")
        assert response.status_code == 200
        with open("page.html", "w") as f:
            f.write(response.data.decode("utf-8"))
    finally:
        if os.path.exists("state.json"):
            os.remove("state.json")

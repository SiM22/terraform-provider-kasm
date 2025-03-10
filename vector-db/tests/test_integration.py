import pytest
from fastapi.testclient import TestClient
from main import app

@pytest.fixture
def client():
    return TestClient(app)


def test_query_endpoint(client):
    response = client.post("/query", json={"text": "test query"})
    assert response.status_code == 200
    assert "results" in response.json()

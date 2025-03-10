from fastapi import FastAPI, HTTPException, Request
from fastapi.responses import HTMLResponse
from fastapi.staticfiles import StaticFiles
from fastapi.templating import Jinja2Templates
from pydantic import BaseModel
import chromadb
import json
import os
from typing import List, Dict, Any, Optional

app = FastAPI(title="Code Structure Vector Database")

# Create directories for static files and templates if they don't exist
os.makedirs("static", exist_ok=True)
os.makedirs("templates", exist_ok=True)

# Mount static files directory
app.mount("/static", StaticFiles(directory="static"), name="static")

# Setup templates
templates = Jinja2Templates(directory="templates")

client = chromadb.Client()
collection = client.create_collection(name="codebase")

class Query(BaseModel):
    text: str

class CodeStructure(BaseModel):
    Package: str
    Imports: List[str]
    Functions: List[Dict[str, Any]]
    Types: List[Dict[str, Any]]
    Variables: Optional[List[Dict[str, Any]]] = None
    ApiMethods: Optional[List[Dict[str, Any]]] = None
    Relations: Optional[List[Dict[str, Any]]] = None

class CodeStructureInput(BaseModel):
    structures: List[CodeStructure]
    source: Optional[str] = None

@app.get("/", response_class=HTMLResponse)
def read_root(request: Request):
    return templates.TemplateResponse("index.html", {"request": request})

@app.get("/api/status")
def api_status():
    return {"status": "ok", "message": "Code Structure Vector Database is running"}

@app.post("/api/add")
def add_code_structure(data: CodeStructureInput):
    try:
        # Convert the data to JSON string for storage
        documents = []
        metadatas = []
        ids = []

        # Get the current count to create unique IDs
        current_count = collection.count()

        for i, structure in enumerate(data.structures):
            # Create a document from the structure
            doc = f"Package: {structure.Package}\n"
            doc += f"Imports: {', '.join(structure.Imports)}\n"

            # Add functions
            for func in structure.Functions:
                doc += f"Function: {func.get('Name', '')} {func.get('Signature', '')}\n"
                if func.get('Doc'):
                    doc += f"Doc: {func.get('Doc')}\n"
                if func.get('Calls') and len(func.get('Calls', [])) > 0:
                    doc += f"Calls: {', '.join(func.get('Calls', []))}\n"

            # Add types
            for typ in structure.Types:
                doc += f"Type: {typ.get('Name', '')}\n"
                if 'Fields' in typ and typ['Fields']:
                    doc += f"Fields: {', '.join(typ.get('Fields', []))}\n"
                if 'Methods' in typ and typ['Methods']:
                    for method in typ['Methods']:
                        doc += f"Method: {method.get('Name', '')} {method.get('Signature', '')}\n"

            # Add variables if they exist
            if structure.Variables:
                for var in structure.Variables:
                    doc += f"Variable: {var.get('Name', '')} {var.get('Type', '')}\n"

            # Add API methods if they exist
            if structure.ApiMethods:
                for api in structure.ApiMethods:
                    doc += f"ApiMethod: {api.get('Name', '')}\n"
                    if api.get('Endpoint'):
                        doc += f"Endpoint: {api.get('Endpoint')}\n"
                    if api.get('HttpMethod'):
                        doc += f"HttpMethod: {api.get('HttpMethod')}\n"
                    if api.get('RequestType'):
                        doc += f"RequestType: {api.get('RequestType')}\n"
                    if api.get('ResponseType'):
                        doc += f"ResponseType: {api.get('ResponseType')}\n"

            # Add relations if they exist
            if structure.Relations:
                for relation in structure.Relations:
                    doc += f"Relation: {relation.get('Source', '')} {relation.get('Type', '')} {relation.get('Target', '')}\n"

            # Create metadata with package and source information
            metadata = {
                "package": structure.Package,
                "source": data.source or "unknown"
            }

            documents.append(doc)
            metadatas.append(metadata)
            ids.append(f"doc_{current_count + i}")

        # Add to collection
        collection.add(
            documents=documents,
            metadatas=metadatas,
            ids=ids
        )

        return {"status": "success", "message": f"Added {len(documents)} code structures to the database"}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/api/query")
def query_codebase(query: Query):
    try:
        results = collection.query(
            query_texts=[query.text],
            n_results=10,
            include=["metadatas", "documents", "distances"]
        )
        return {"results": results}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/api/clear")
def clear_database():
    try:
        # Delete the collection and recreate it
        client.delete_collection(name="codebase")
        global collection
        collection = client.create_collection(name="codebase")
        return {"status": "success", "message": "Database cleared successfully"}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/api/count")
def get_count():
    try:
        count = collection.count()
        return {"count": count}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/api/all")
def get_all_documents():
    try:
        # Get all documents from the collection
        results = collection.get(
            include=["metadatas", "documents"]
        )
        return {"results": results}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/api/flow")
def get_code_flow():
    try:
        # Get all documents from the collection
        results = collection.get(
            include=["metadatas", "documents"]
        )

        # Extract all relations from the documents
        nodes = set()
        edges = []

        for doc in results["documents"]:
            lines = doc.split('\n')
            for line in lines:
                if line.startswith("Relation:"):
                    parts = line.replace("Relation: ", "").split()
                    if len(parts) >= 3:
                        source = parts[0]
                        rel_type = parts[1]
                        target = parts[2]

                        nodes.add(source)
                        nodes.add(target)

                        edges.append({
                            "source": source,
                            "target": target,
                            "type": rel_type
                        })

        return {
            "nodes": list(nodes),
            "edges": edges
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/api/api_methods")
def get_api_methods():
    try:
        # Get all documents from the collection
        results = collection.get(
            include=["metadatas", "documents"]
        )

        # Extract all API methods from the documents
        api_methods = []

        for doc in results["documents"]:
            lines = doc.split('\n')
            current_api = None

            for line in lines:
                if line.startswith("ApiMethod:"):
                    if current_api:
                        api_methods.append(current_api)

                    name = line.replace("ApiMethod: ", "")
                    current_api = {"name": name}
                elif line.startswith("Endpoint:") and current_api:
                    current_api["endpoint"] = line.replace("Endpoint: ", "")
                elif line.startswith("HttpMethod:") and current_api:
                    current_api["httpMethod"] = line.replace("HttpMethod: ", "")
                elif line.startswith("RequestType:") and current_api:
                    current_api["requestType"] = line.replace("RequestType: ", "")
                elif line.startswith("ResponseType:") and current_api:
                    current_api["responseType"] = line.replace("ResponseType: ", "")

            if current_api:
                api_methods.append(current_api)

        return {"api_methods": api_methods}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/api/sources")
def get_sources():
    try:
        # Get all documents with their metadata
        results = collection.get(
            include=["metadatas"]
        )

        # Extract unique sources from metadata
        sources = set()
        if results and "metadatas" in results and results["metadatas"]:
            for metadata in results["metadatas"]:
                if metadata and "source" in metadata:
                    sources.add(metadata["source"])

        return {"sources": list(sources)}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)

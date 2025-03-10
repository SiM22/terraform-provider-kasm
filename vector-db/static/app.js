document.addEventListener('DOMContentLoaded', function() {
    // DOM elements
    const searchInput = document.getElementById('searchInput');
    const searchButton = document.getElementById('searchButton');
    const navigationMenu = document.getElementById('navigationMenu');
    const contentTitle = document.getElementById('contentTitle');
    const contentArea = document.getElementById('contentArea');
    const documentCount = document.getElementById('documentCount');
    const loadingIndicator = document.getElementById('loadingIndicator');
    const detailCard = document.getElementById('detailCard');
    const detailTitle = document.getElementById('detailTitle');
    const detailContent = document.getElementById('detailContent');

    // Load D3.js for visualizations
    const d3Script = document.createElement('script');
    d3Script.src = 'https://d3js.org/d3.v7.min.js';
    document.head.appendChild(d3Script);

    // Global state
    let allDocuments = [];
    let currentView = 'packages';
    let codeFlowData = null;
    let apiMethodsData = null;

    // Initialize the app
    init();

    // Event listeners
    searchButton.addEventListener('click', performSearch);
    searchInput.addEventListener('keyup', function(event) {
        if (event.key === 'Enter') {
            performSearch();
        }
    });

    navigationMenu.addEventListener('click', function(event) {
        if (event.target.tagName === 'A') {
            event.preventDefault();

            // Update active link
            document.querySelectorAll('#navigationMenu .nav-link').forEach(link => {
                link.classList.remove('active');
            });
            event.target.classList.add('active');

            // Update view
            currentView = event.target.getAttribute('data-view');

            // For code flow and API methods views, fetch the data if not already loaded
            if (currentView === 'code-flow' && !codeFlowData) {
                fetchCodeFlowData();
            } else if (currentView === 'api-methods' && !apiMethodsData) {
                fetchApiMethodsData();
            } else {
                renderContent();
            }
        }
    });

    // Functions
    async function init() {
        showLoading(true);
        try {
            // Get document count
            const countResponse = await fetch('/api/count');
            const countData = await countResponse.json();
            documentCount.textContent = `${countData.count} documents`;

            // Get all documents
            const response = await fetch('/api/all');
            const data = await response.json();
            allDocuments = data.results;

            // Update navigation menu to include new views
            const navMenu = document.getElementById('navigationMenu');
            navMenu.innerHTML += `
                <li class="nav-item">
                    <a class="nav-link" href="#" data-view="code-flow">Code Flow</a>
                </li>
                <li class="nav-item">
                    <a class="nav-link" href="#" data-view="api-methods">API Methods</a>
                </li>
            `;

            renderContent();
        } catch (error) {
            console.error('Error initializing app:', error);
            contentArea.innerHTML = `<div class="alert alert-danger">Error loading data: ${error.message}</div>`;
        } finally {
            showLoading(false);
        }
    }

    async function fetchCodeFlowData() {
        showLoading(true);
        try {
            const response = await fetch('/api/flow');
            if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
            }
            codeFlowData = await response.json();
            console.log('Code flow data:', codeFlowData);
            renderContent();
        } catch (error) {
            console.error('Error fetching code flow data:', error);
            contentArea.innerHTML = `<div class="alert alert-danger">Error loading code flow data: ${error.message}</div>`;
        } finally {
            showLoading(false);
        }
    }

    async function fetchApiMethodsData() {
        showLoading(true);
        try {
            const response = await fetch('/api/api_methods');
            if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
            }
            apiMethodsData = await response.json();
            console.log('API methods data:', apiMethodsData);
            renderContent();
        } catch (error) {
            console.error('Error fetching API methods data:', error);
            contentArea.innerHTML = `<div class="alert alert-danger">Error loading API methods data: ${error.message}</div>`;
        } finally {
            showLoading(false);
        }
    }

    function showLoading(show) {
        loadingIndicator.style.display = show ? 'block' : 'none';
        contentArea.style.display = show ? 'none' : 'block';
    }

    function renderContent() {
        // Hide detail card when changing views
        detailCard.style.display = 'none';

        if (currentView === 'codeFlow') {
            renderCodeFlow();
            return;
        } else if (currentView === 'apiMethods') {
            renderApiMethods();
            return;
        }

        if (!allDocuments || !allDocuments.documents) {
            contentArea.innerHTML = '<div class="alert alert-info">No documents found in the database.</div>';
            return;
        }

        // Set content title based on current view
        contentTitle.textContent = currentView.charAt(0).toUpperCase() + currentView.slice(1);

        let html = '';

        if (currentView === 'packages') {
            // Group by package
            const packages = new Set();
            if (allDocuments.metadatas) {
                allDocuments.metadatas.forEach(metadata => {
                    if (metadata.package) {
                        packages.add(metadata.package);
                    }
                });
            }

            html = '<div class="row">';
            Array.from(packages).forEach(packageName => {
                html += `
                    <div class="col-md-6">
                        <div class="card item-card" onclick="showPackageDetails('${packageName}')">
                            <div class="card-body">
                                <h6 class="package-name">${packageName}</h6>
                                <p>Package</p>
                            </div>
                        </div>
                    </div>
                `;
            });
            html += '</div>';
        } else if (currentView === 'functions' || currentView === 'types') {
            // Extract functions or types from documents
            const items = new Map();

            allDocuments.documents.forEach((doc, index) => {
                const lines = doc.split('\n');
                const packageLine = lines.find(line => line.startsWith('Package:'));
                let packageName = 'Unknown';

                if (packageLine) {
                    packageName = packageLine.replace('Package:', '').trim();
                }

                lines.forEach(line => {
                    if (currentView === 'functions' && line.startsWith('Function:')) {
                        const functionName = line.replace('Function:', '').trim();
                        if (functionName) {
                            items.set(functionName, {
                                name: functionName,
                                package: packageName,
                                docId: allDocuments.ids[index]
                            });
                        }
                    } else if (currentView === 'types' && line.startsWith('Type:')) {
                        const typeName = line.replace('Type:', '').trim();
                        if (typeName) {
                            items.set(typeName, {
                                name: typeName,
                                package: packageName,
                                docId: allDocuments.ids[index]
                            });
                        }
                    }
                });
            });

            html = '<div class="row">';
            Array.from(items.values()).forEach(item => {
                const itemClass = currentView === 'functions' ? 'function-name' : 'type-name';
                html += `
                    <div class="col-md-6">
                        <div class="card item-card" onclick="showItemDetails('${item.docId}', '${item.name}', '${currentView}')">
                            <div class="card-body">
                                <h6 class="${itemClass}">${item.name}</h6>
                                <p>Package: ${item.package}</p>
                            </div>
                        </div>
                    </div>
                `;
            });
            html += '</div>';
        }

        contentArea.innerHTML = html;

        // Add global functions for onclick handlers
        window.showPackageDetails = function(packageName) {
            showDetails(packageName, 'package');
        };

        window.showItemDetails = function(docId, itemName, itemType) {
            showDetails(docId, itemType, itemName);
        };
    }

    function showDetails(id, type, name = null) {
        let content = '';
        let title = '';

        if (type === 'package') {
            // Find all documents for this package
            const packageDocs = [];
            allDocuments.metadatas.forEach((metadata, index) => {
                if (metadata.package === id) {
                    packageDocs.push({
                        document: allDocuments.documents[index],
                        id: allDocuments.ids[index]
                    });
                }
            });

            title = `Package: ${id}`;
            content = packageDocs.map(doc => formatDocument(doc.document)).join('\n\n');
        } else {
            // Find the specific document
            const index = allDocuments.ids.indexOf(id);
            if (index !== -1) {
                const document = allDocuments.documents[index];

                if (type === 'functions') {
                    title = `Function: ${name}`;
                    // Extract function details
                    const lines = document.split('\n');
                    const functionLines = [];
                    let inFunction = false;

                    for (const line of lines) {
                        if (line.startsWith('Function:') && line.includes(name)) {
                            inFunction = true;
                            functionLines.push(line);
                        } else if (inFunction && line.startsWith('Doc:')) {
                            functionLines.push(line);
                        } else if (inFunction && (line.startsWith('Function:') || line.startsWith('Type:'))) {
                            inFunction = false;
                        }
                    }

                    content = functionLines.join('\n');
                } else if (type === 'types') {
                    title = `Type: ${name}`;
                    // Extract type details
                    const lines = document.split('\n');
                    const typeLine = lines.find(line => line.startsWith('Type:') && line.includes(name));
                    content = typeLine || 'Type details not found';
                }
            } else {
                content = 'Document not found';
            }
        }

        detailTitle.textContent = title;
        detailContent.textContent = content;
        detailCard.style.display = 'block';
    }

    function formatDocument(doc) {
        return doc.split('\n').map(line => {
            if (line.startsWith('Package:')) {
                return `<span class="package-name">${line}</span>`;
            } else if (line.startsWith('Function:')) {
                return `<span class="function-name">${line}</span>`;
            } else if (line.startsWith('Type:')) {
                return `<span class="type-name">${line}</span>`;
            } else if (line.startsWith('Imports:')) {
                return `<span class="import-item">${line}</span>`;
            } else if (line.startsWith('Doc:')) {
                return `<span class="doc-text">${line}</span>`;
            }
            return line;
        }).join('\n');
    }

    function renderCodeFlow() {
        contentTitle.textContent = 'Code Flow Visualization';

        if (!codeFlowData || !codeFlowData.nodes || !codeFlowData.edges || codeFlowData.nodes.length === 0) {
            contentArea.innerHTML = '<div class="alert alert-info">No code flow data available.</div>';
            return;
        }

        // Create a container for the visualization with controls
        contentArea.innerHTML = `
            <div class="card mb-3">
                <div class="card-body">
                    <div class="row">
                        <div class="col-md-6">
                            <div class="form-group">
                                <label for="nodeSearch">Search Nodes:</label>
                                <input type="text" id="nodeSearch" class="form-control" placeholder="Type to search...">
                            </div>
                        </div>
                        <div class="col-md-6">
                            <div class="form-group">
                                <label>Node Categories:</label>
                                <div id="nodeCategories" class="mt-2"></div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
            <div id="flow-container" style="height: 600px; border: 1px solid #ddd; border-radius: 5px;"></div>
            <div id="node-details" class="mt-3 d-none">
                <div class="card">
                    <div class="card-header bg-primary text-white">
                        <h5 class="mb-0" id="node-details-title">Node Details</h5>
                    </div>
                    <div class="card-body" id="node-details-content"></div>
                </div>
            </div>
        `;

        // Wait for D3.js to load
        if (typeof d3 === 'undefined') {
            setTimeout(renderCodeFlow, 100);
            return;
        }

        const container = document.getElementById('flow-container');
        const width = container.clientWidth;
        const height = container.clientHeight;
        const nodeSearch = document.getElementById('nodeSearch');
        const nodeDetailsPanel = document.getElementById('node-details');
        const nodeDetailsTitle = document.getElementById('node-details-title');
        const nodeDetailsContent = document.getElementById('node-details-content');
        const nodeCategoriesContainer = document.getElementById('nodeCategories');

        // Process nodes and edges to add more information
        const nodes = codeFlowData.nodes.map(d => {
            // Determine node type based on naming conventions
            let nodeType = 'function';
            if (d.includes('Client')) nodeType = 'client';
            else if (d.includes('Resource')) nodeType = 'resource';
            else if (d.includes('DataSource')) nodeType = 'datasource';
            else if (d.includes('Provider')) nodeType = 'provider';
            else if (d.includes('API') || d.includes('Api')) nodeType = 'api';

            return {
                id: d,
                name: d,
                type: nodeType,
                connections: 0 // Will count connections later
            };
        });

        // Count connections for each node
        const edges = codeFlowData.edges.map(d => {
            // Find source and target nodes
            const sourceNode = nodes.find(n => n.id === d.source);
            const targetNode = nodes.find(n => n.id === d.target);

            // Increment connection count
            if (sourceNode) sourceNode.connections++;
            if (targetNode) targetNode.connections++;

            return {
                source: d.source,
                target: d.target,
                type: d.type || 'calls'
            };
        });

        // Create color scale for node types
        const nodeTypes = [...new Set(nodes.map(n => n.type))];
        const colorScale = d3.scaleOrdinal()
            .domain(nodeTypes)
            .range(['#69b3a2', '#ff7f0e', '#1f77b4', '#d62728', '#9467bd', '#8c564b']);

        // Add node type filters
        nodeCategoriesContainer.innerHTML = nodeTypes.map(type => {
            const color = colorScale(type);
            return `
                <div class="form-check form-check-inline">
                    <input class="form-check-input" type="checkbox" id="filter-${type}" value="${type}" checked>
                    <label class="form-check-label" for="filter-${type}">
                        <span class="badge" style="background-color: ${color}">${type}</span>
                    </label>
                </div>
            `;
        }).join('');

        // Create the SVG container
        const svg = d3.select('#flow-container')
            .append('svg')
            .attr('width', width)
            .attr('height', height);

        // Create a group for the graph
        const g = svg.append('g');

        // Add zoom behavior
        const zoom = d3.zoom()
            .scaleExtent([0.1, 4])
            .on('zoom', (event) => {
                g.attr('transform', event.transform);
            });

        svg.call(zoom);

        // Create the force simulation
        const simulation = d3.forceSimulation(nodes)
            .force('link', d3.forceLink(edges)
                .id(d => d.id)
                .distance(d => 100 + (d.source.connections + d.target.connections) * 2))
            .force('charge', d3.forceManyBody().strength(-300))
            .force('center', d3.forceCenter(width / 2, height / 2))
            .force('collision', d3.forceCollide().radius(d => 10 + Math.min(d.connections * 2, 20)));

        // Create the links
        const link = g.selectAll('.link')
            .data(edges)
            .enter().append('line')
            .attr('class', 'link')
            .attr('stroke', d => d.type === 'implements' ? '#ff7f0e' : '#999')
            .attr('stroke-opacity', 0.6)
            .attr('stroke-width', d => 1 + Math.min(d.source.connections, 3))
            .attr('marker-end', d => `url(#arrow-${d.type})`);

        // Add arrow markers for different link types
        const markerTypes = [...new Set(edges.map(e => e.type))];
        const defs = svg.append('defs');

        markerTypes.forEach(type => {
            defs.append('marker')
                .attr('id', `arrow-${type}`)
                .attr('viewBox', '0 -5 10 10')
                .attr('refX', 20)
                .attr('refY', 0)
                .attr('markerWidth', 6)
                .attr('markerHeight', 6)
                .attr('orient', 'auto')
                .append('path')
                .attr('d', 'M0,-5L10,0L0,5')
                .attr('fill', type === 'implements' ? '#ff7f0e' : '#999');
        });

        // Create the nodes
        const node = g.selectAll('.node')
            .data(nodes)
            .enter().append('g')
            .attr('class', 'node')
            .attr('data-id', d => d.id)
            .attr('data-type', d => d.type)
            .call(d3.drag()
                .on('start', dragstarted)
                .on('drag', dragged)
                .on('end', dragended))
            .on('click', showNodeDetails);

        // Add circles for the nodes
        node.append('circle')
            .attr('r', d => 8 + Math.min(d.connections, 10))
            .attr('fill', d => colorScale(d.type))
            .attr('stroke', '#fff')
            .attr('stroke-width', 1.5);

        // Add labels to the nodes
        node.append('text')
            .attr('dx', d => 10 + Math.min(d.connections, 10))
            .attr('dy', '.35em')
            .text(d => d.id.length > 20 ? d.id.substring(0, 20) + '...' : d.id)
            .attr('font-size', '10px');

        // Update positions on each tick of the simulation
        simulation.on('tick', () => {
            link
                .attr('x1', d => d.source.x)
                .attr('y1', d => d.source.y)
                .attr('x2', d => d.target.x)
                .attr('y2', d => d.target.y);

            node
                .attr('transform', d => `translate(${d.x},${d.y})`);
        });

        // Search functionality
        nodeSearch.addEventListener('input', function() {
            const searchTerm = this.value.toLowerCase();

            // Highlight matching nodes
            node.each(function(d) {
                const nodeElement = d3.select(this);
                const matches = d.id.toLowerCase().includes(searchTerm);

                nodeElement.select('circle')
                    .transition()
                    .duration(300)
                    .attr('r', matches && searchTerm ? 12 + Math.min(d.connections, 10) : 8 + Math.min(d.connections, 10))
                    .attr('fill', matches && searchTerm ? '#3498db' : colorScale(d.type));

                nodeElement.select('text')
                    .transition()
                    .duration(300)
                    .attr('font-weight', matches && searchTerm ? 'bold' : 'normal');
            });
        });

        // Node type filtering
        nodeCategoriesContainer.addEventListener('change', function(e) {
            if (e.target.type === 'checkbox') {
                const nodeType = e.target.value;
                const isChecked = e.target.checked;

                // Filter nodes
                node.filter(d => d.type === nodeType)
                    .style('display', isChecked ? 'block' : 'none');

                // Filter links connected to these nodes
                link.style('display', d => {
                    const sourceVisible = node.filter(n => n.id === d.source.id).style('display') !== 'none';
                    const targetVisible = node.filter(n => n.id === d.target.id).style('display') !== 'none';
                    return sourceVisible && targetVisible ? 'block' : 'none';
                });
            }
        });

        // Show node details
        function showNodeDetails(event, d) {
            // Find incoming and outgoing connections
            const incoming = edges.filter(e => e.target.id === d.id);
            const outgoing = edges.filter(e => e.source.id === d.id);

            nodeDetailsTitle.textContent = d.id;
            nodeDetailsContent.innerHTML = `
                <p><strong>Type:</strong> <span class="badge" style="background-color: ${colorScale(d.type)}">${d.type}</span></p>
                <p><strong>Connections:</strong> ${d.connections}</p>

                <div class="row mt-3">
                    <div class="col-md-6">
                        <h6>Incoming Connections (${incoming.length})</h6>
                        <ul class="list-group">
                            ${incoming.length ? incoming.map(e => `
                                <li class="list-group-item">
                                    <span class="badge bg-secondary">${e.type}</span>
                                    ${e.source.id}
                                </li>
                            `).join('') : '<li class="list-group-item">None</li>'}
                        </ul>
                    </div>
                    <div class="col-md-6">
                        <h6>Outgoing Connections (${outgoing.length})</h6>
                        <ul class="list-group">
                            ${outgoing.length ? outgoing.map(e => `
                                <li class="list-group-item">
                                    <span class="badge bg-secondary">${e.type}</span>
                                    ${e.target.id}
                                </li>
                            `).join('') : '<li class="list-group-item">None</li>'}
                        </ul>
                    </div>
                </div>
            `;

            nodeDetailsPanel.classList.remove('d-none');
        }

        // Functions for drag behavior
        function dragstarted(event, d) {
            if (!event.active) simulation.alphaTarget(0.3).restart();
            d.fx = d.x;
            d.fy = d.y;
        }

        function dragged(event, d) {
            d.fx = event.x;
            d.fy = event.y;
        }

        function dragended(event, d) {
            if (!event.active) simulation.alphaTarget(0);
            // Keep the node fixed where the user dragged it
            // d.fx = null;
            // d.fy = null;
        }
    }

    function renderApiMethods() {
        contentTitle.textContent = 'API Methods';

        if (!apiMethodsData || !apiMethodsData.api_methods || apiMethodsData.api_methods.length === 0) {
            contentArea.innerHTML = '<div class="alert alert-info">No API methods data available.</div>';
            return;
        }

        let html = '<div class="row">';

        apiMethodsData.api_methods.forEach(method => {
            const endpoint = method.endpoint || 'Unknown';
            const httpMethod = method.httpMethod ? method.httpMethod.toLowerCase() : 'unknown';
            const requestType = method.requestType || 'None';
            const responseType = method.responseType || 'None';

            // Determine HTTP method class
            const methodClass = ['get', 'post', 'put', 'delete', 'patch'].includes(httpMethod) ? httpMethod : 'get';

            html += `
                <div class="col-md-6 mb-3">
                    <div class="card api-method-card">
                        <div class="card-header bg-primary text-white">
                            <h5 class="mb-0">${method.name}</h5>
                        </div>
                        <div class="card-body">
                            <div>
                                <span class="http-method ${methodClass}">${httpMethod.toUpperCase()}</span>
                                <span class="api-endpoint">${endpoint}</span>
                            </div>
                            <div class="request-response-container">
                                <div class="row">
                                    <div class="col-md-6">
                                        <div class="card request-response-card">
                                            <div class="card-header bg-light">Request</div>
                                            <div class="card-body">
                                                <code>${requestType}</code>
                                            </div>
                                        </div>
                                    </div>
                                    <div class="col-md-6">
                                        <div class="card request-response-card">
                                            <div class="card-header bg-light">Response</div>
                                            <div class="card-body">
                                                <code>${responseType}</code>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            `;
        });

        html += '</div>';
        contentArea.innerHTML = html;
    }

    async function performSearch() {
        const query = searchInput.value.trim();
        if (!query) return;

        showLoading(true);
        try {
            const response = await fetch('/api/query', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ text: query })
            });

            const data = await response.json();

            contentTitle.textContent = `Search Results: "${query}"`;

            if (data.results && data.results.documents && data.results.documents[0].length > 0) {
                let html = '<div class="alert alert-success">Found ' + data.results.documents[0].length + ' results</div>';

                html += '<div class="list-group">';
                data.results.documents[0].forEach((doc, index) => {
                    const distance = data.results.distances[0][index];
                    const relevance = Math.round((1 - distance) * 100);

                    html += `
                        <div class="list-group-item">
                            <div class="d-flex justify-content-between align-items-center">
                                <h6>Result #${index + 1}</h6>
                                <span class="badge bg-primary">${relevance}% match</span>
                            </div>
                            <pre class="code-block">${formatDocument(doc)}</pre>
                        </div>
                    `;
                });
                html += '</div>';

                contentArea.innerHTML = html;
            } else {
                contentArea.innerHTML = '<div class="alert alert-warning">No results found for your query.</div>';
            }
        } catch (error) {
            console.error('Error searching:', error);
            contentArea.innerHTML = `<div class="alert alert-danger">Error performing search: ${error.message}</div>`;
        } finally {
            showLoading(false);
        }
    }
});

package az

const documentationTemplate = `<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Name}}</title>
		<style>
			html {
				font-family: monospace;
			}
			.endpoint {
				border: solid 1px #c5c3c3;
				padding: 10px;
				margin-bottom: 10px;
				background-color: #e4f9ff;
			}
			.MethodStruct {
				background-color: #CDDC39;
				border-radius: 5px;
				padding: 3px;
				font-weight: bold;
			}
			.GET {
				background-color: #4caf50;
			}
			.POST {
				background-color: #03a9f4;
			}
			.DELETE {
				background-color: #ff5722;
			}
			.PUT {
				background-color: #74dfed;
			}
			.OPTIONS {
				background-color: #ffc107;
			}
			.HEAD {
				background-color: #b5bcc0;
			}

			table {
				display: block;
				width: 100%;
				overflow: auto;
				margin-top: 0;
				margin-bottom: 16px;
				border-spacing: 0;
  				border-collapse: collapse;
			}

			table th {
				font-weight: 600;
			}

			table th,
			table td {
				padding: 6px 13px;
				border: 1px solid #dfe2e5;
			}

			table tr {
				background-color: #fff;
				border-top: 1px solid #c6cbd1;
			}

			table tr:nth-child(2n) {
				background-color: #f6f8fa;
			}


			#response {
				display: none;
				position: fixed; 
				z-index: 1;
				left: 0;
				top: 0;
				width: 100%;
				height: 100%;
				overflow: auto; 
				background-color: rgb(0,0,0);
				background-color: rgba(0,0,0,0.4);
			}

			#response-content {
				background-color: #fefefe;
				margin: 15% auto; 
				padding: 20px;
				border: 1px solid #888;
				width: 80%;
			}

			#response-close {
				color: #aaa;
				float: right;
				font-size: 28px;
				font-weight: bold;
			}

			#response-close:hover,
			#response-close:focus {
				color: black;
				text-decoration: none;
				cursor: pointer;
			}

			.api-test {
				background-color: #eafdd5;
				border: solid 1px #c5c3c3;
				padding: 10px;
			}

			input, select, textarea {
				margin-bottom: 5px;
			}
			input[type="submit"] {
				margin: 0;
			}

			input[type="submit"] {
				margin: 0;
			}
			input[type="text"] {
				width: 100%;
			}

			.description {
				white-space:pre;
			}


		</style>
	</head>
	<body>
	<h1>{{.Name}}</h1>
	<p class="description">{{.Description}}</p>
	<hr>
	<button onclick='
		for(var i =0, il = document.getElementsByClassName("collapse").length;i<il;i++){
			document.getElementsByClassName("collapse")[i].style.display = "block";
		}
		for(var i =0, il = document.getElementsByClassName("collapse-toogle").length;i<il;i++){
			document.getElementsByClassName("collapse-toogle")[i].innerText = "Collapse";
		}
	'>Expand All</button>

	<button  onclick='
		for(var i =0, il = document.getElementsByClassName("collapse").length;i<il;i++){
			document.getElementsByClassName("collapse")[i].style.display = "none";
		}
		for(var i =0, il = document.getElementsByClassName("collapse-toogle").length;i<il;i++){
			document.getElementsByClassName("collapse-toogle")[i].innerText = "Expand";
		}
	'>Collapse All</button>
	{{ range $key, $value := .Namespace }}
		<h2>{{ if $key }}{{ $key }}{{ else }}WITHOUT NAMESPACE{{ end }} <button class="collapse-toogle" style="cursor: pointer;" onclick="toogle(event)">Expand</button></h2>
			<div class="collapse" style="display: none;">
			{{range $value}}
				<div class="endpoint">
					<div class="path-MethodStruct">
						<span class="MethodStruct {{ .Method }}">{{ if .Method }}{{ .Method }}{{ else }}ANY{{ end }}</span>
						<span class="path"><strong>{{ .Namespace }}</strong>{{replace .Path .Namespace ""}}{{ if .RequiredParams}}{{ range $i, $p := .RequiredParams }}{{if eq $i 0}}?{{ $p.Name }}={{"{"}}:{{ $p.Name }}{{"}"}}{{ else }}&{{ $p.Name }}={{"{"}}:{{ $p.Name }}{{"}"}}{{ end }}{{ end }}{{ end }}</span>
					</div>
					<hr>
					<h3 style="margin: 0;" >{{- .Name -}} <button  class="collapse-toogle" style="margin-top: -35px; float: right; cursor: pointer;" onclick="toogle(event)">Expand</button></h3>
					<div class="collapse" style="display: none;">
						<p class="description">{{.Description}}</p>
						

						{{ if .RequiredParams}} 
						<strong>Required Parameters</strong>
						<table style="width:100%">
							<tr>
								<th>Parameter</th>
								<th>Type</th> 
								<th>Description</th>
							</tr>
							{{ range .RequiredParams }}
							<tr>
								<td>{{- .Name -}}</td>
								<td>{{- .ParamType -}}</td> 
								<td class="description">{{.Description}}</td>
							</tr>
							{{ end }}
						</table>
						{{ end }}
						
						{{ if .DocParams}} 
						<strong>Optional Parameters</strong>
						<table style="width:100%">
							<tr>
								<th>Parameter</th>
								<th>Type</th> 
								<th>Description</th>
							</tr>
							{{ range .DocParams }}
							<tr>
								<td>{{- .Name -}}</td>
								<td>{{- .ParamType -}}</td> 
								<td class="description">{{ .Description }}</td>
							</tr>
							{{ end }}
						</table>
						{{ end }}

						{{ if .DocHeaders}} 
						<strong>Headers</strong>
						<table style="width:100%">
							<tr>
								<th>Header</th>
								<th>Description</th>
							</tr>
							{{ range .DocHeaders }}
							<tr>
								<td>{{- .Name -}}</td>
								<td class="description">{{ .Description }}</td>
							</tr>
							{{ end }}
						</table>
						{{ end }}

						<form class="api-test" name="{{.Method}}" action="{{.Path}}" onsubmit="APITest(event)">
							<h3>Test API</h3>
							{{ if not .Method }}
								Method<br>
								<select onchange='changeMethod(this)'>
									<option value="GET">GET</option>
									<option value="POST">POST</option>
									<option value="PUT">PUT</option>
									<option value="DELETE">DELETE</option>
									<option value="OPTIONS">OPTIONS</option>
									<option value="HEAD">HEAD</option>
								</select><br>
							{{ end }}
							{{ if .RequiredParams}} 
								<h4>Required Parameters</h4>
								{{ range .RequiredParams }}
									{{- .Name -}}<br>
									<input type="text" name="{{ .Name }}" alt="req-ParamStruct"><br>
								{{ end }}
							{{ end }}

							{{ if .DocParams}} 
								<h4>Optional Parameters</h4>
								{{ range .DocParams }}
									{{- .Name -}}<br>
									<input type="text" name="{{ .Name }}" alt="ParamStruct"><br>
								{{ end }}
							{{ end }}

							{{ if .DocHeaders}} 
							<h4>Headers</h4>
								{{ range .DocHeaders }}
									{{- .Name -}}<br>
									<input type="text" name="{{ .Name }}" alt="header" list="{{ .Name | ToLower}}"><br>
								{{ end }}
							{{ end }}
							{{if ne .Method "GET"}}
								{{if ne .Method "HEAD"}}
									<h4 class="body-label" {{ if not .Method }}style="display: none;"{{ end }}>Body</h4>
									<textarea style="width:100%;{{ if not .Method }}display: none;{{ end }}" rows="10" name="body"></textarea>
								{{end}}
							{{end}}
							<br>
							<strong>Extra Parameters</strong><br>
							<input type="text" alt="extra-params"><br>
							<input type="submit" value="Send">
						</form>
					</div>
				</div>
			{{end}}
		</div>
	{{ end }}

	<div id="response">
		<div id="response-content">
			<span id="response-close">&times;</span>
			<strong><span id="request-MethodStruct"></span></strong> <span id="request-url"></span><br>
			<pre id="request-body" style="overflow: auto;"></pre>
			<hr>
			<strong>Date/Time:</strong> <span id="response-date"></span><br>
			<strong>Content-Type:</strong> <span id="response-content-type"></span><br>
			<strong>Status:</strong> <span id="response-status"></span> <span id="response-status-text"></span>
			<pre id="response-body" style="overflow: auto;"></pre>
		</div>
	</div>


	<datalist id="content-type">
		<option value="text/plain">
		<option value="application/json">
		<option value="text/html">
		<option value="application/octet-stream">
		<option value="text/css">
		<option value="text/csv">
		<option value="image/gif">
		<option value="image/x-icon">
		<option value="text/calendar">
		<option value="image/jpeg">
		<option value="application/javascript">
		<option value="video/mpeg">
		<option value="audio/ogg">
		<option value="video/ogg">
		<option value="application/ogg">
		<option value="font/otf">
		<option value="image/png">
		<option value="application/pdf">
		<option value="image/svg+xml">
		<option value="image/tiff">
		<option value="application/typescript">
		<option value="font/ttf">
		<option value="audio/x-wav">
		<option value="audio/webm">
		<option value="video/webm">
		<option value="font/woff">
		<option value="font/woff2">
		<option value="application/xhtml+xml">
		<option value="application/xml">
		<option value="application/zip">
		<option value="video/3gpp">
	</datalist>




		<script>


			function changeMethod(m) {
				m.parentElement.name = m.options[m.selectedIndex].value;
				if (m.options[m.selectedIndex].value == "GET" || m.options[m.selectedIndex].value == "HEAD")
					{
						m.parentElement.getElementsByTagName("textarea")[0].style.display = "none";
						m.parentElement.getElementsByClassName("body-label")[0].style.display = "none";
					} else {
						m.parentElement.getElementsByTagName("textarea")[0].style.display = "block";
						m.parentElement.getElementsByClassName("body-label")[0].style.display = "block";
					}
			}



	
			var response = document.getElementById('response');

			var close = document.getElementById("response-close");

			close.onclick = function() {
				response.style.display = "none";
			}

			window.onclick = function(event) {
				if (event.target == response) {
					response.style.display = "none";
				}
			}



			function APITest(e) {
				e.preventDefault()

				var params = ""
				var extraParams = ""
				var body = ""
				var headers = {}
				for (i = 0; i < e.target.length; i++) {
					if (e.target[i].tagName == "INPUT" && e.target[i].alt == "req-ParamStruct") {
						if (params == "") {
							params = params+"?"+e.target[i].name+"="+encodeURIComponent(e.target[i].value)
						} else {
							params = params+"&"+e.target[i].name+"="+encodeURIComponent(e.target[i].value)
						}
					}
					if (e.target[i].tagName == "INPUT" && e.target[i].alt == "ParamStruct" && e.target[i].value != "") {
						if (params == "") {
							params = params+"?"+e.target[i].name+"="+encodeURIComponent(e.target[i].value)
						} else {
							params = params+"&"+e.target[i].name+"="+encodeURIComponent(e.target[i].value)
						}
					}
					if (e.target[i].tagName == "INPUT" && e.target[i].alt == "header" && e.target[i].value != "") {
						headers[e.target[i].name] = e.target[i].value
					}
					if (e.target[i].tagName == "INPUT" && e.target[i].alt == "extra-params") {
						extraParams = e.target[i].value
					}
					if (e.target[i].tagName == "TEXTAREA" && e.target[i].name == "body") {
						body = e.target[i].value
					}
				}

				MethodStruct = e.target.attributes.name.nodeValue || "GET"

				option = {
					method: MethodStruct
				}

				if (headers) {
					option["headers"] = headers
				}

				if (body && MethodStruct != "GET" && MethodStruct != "HEAD") {
					option["body"] = body
				}

				if (extraParams) {
					if (params) {
						extraParams = "&"+extraParams
					} else {
						extraParams = "?"+extraParams
					}
				}

				url = e.target.action+params+extraParams
				
				fetch(url, option).then(function(res) {
					contentType = res.headers.get('Content-Type');
					date = res.headers.get('Date')
					status = res.status
					statusText = res.statusText
					return res.text()
				}).then(function(b) {
					document.getElementById("request-MethodStruct").innerText = MethodStruct;
					document.getElementById("request-url").innerText = url;
					if (body) {
						try {
							bodyJSON = JSON.parse(body)
							body = JSON.stringify(bodyJSON, null, 2)
						} catch(e) {}
					}
					if (MethodStruct == "GET" || MethodStruct == "HEAD"){
						document.getElementById("request-body").innerText = "";
					} else {
						document.getElementById("request-body").innerText = body;
					}
					document.getElementById("response-date").innerText = date;
					document.getElementById("response-content-type").innerText = contentType;
					document.getElementById("response-status").innerText = status;
					document.getElementById("response-status-text").innerText = statusText;
					if (b) {
						try {
							bodyJSON = JSON.parse(b)
							b = JSON.stringify(bodyJSON, null, 2)
						} catch(e) {
							if (contentType == "application/json") {
								alert(e);
							}
						}
					}
					document.getElementById("response-body").innerText = b;
					document.getElementById('response').style.display = "block";
				})
			}



			function toogle(e) {
				if (e.target.parentElement.nextElementSibling.style.display == "none") {
					e.target.parentElement.nextElementSibling.style.display = "block"
					e.target.innerText = "Collapse"
				} else {
					e.target.parentElement.nextElementSibling.style.display = "none"
					e.target.innerText = "Expand"
				}
			}


		</script>
	</body>
</html>`

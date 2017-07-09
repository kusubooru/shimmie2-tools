// generated by go generate; DO NOT EDIT

package main

import "html/template"

var (
	layoutTmpl = template.Must(template.New("layout").Funcs(fns).Parse(layoutTemplate))

	indexTmpl = template.Must(template.Must(layoutTmpl.Clone()).Parse(indexTemplate))

	uploadTmpl = template.Must(template.Must(layoutTmpl.Clone()).Parse(uploadTemplate))

	layoutTemplate = `
{{ define "layout" }}
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <title>Tagaa {{printv .Version}}</title>
    <meta name="description" content="Interface for the 'Bulk Add CSV' Shimmie2 extension">
    <meta name="author" content="kusubooru">

    <style>
      html {
        font-family: sans-serif;
      }
      input {
        margin-bottom: 0.6em;
      }
      .block {
        display: block;
        padding: 15px;
        margin-bottom: 10px;
      }
      .block-danger {
        background: #f2dede;
        color: #333;
      }
      .block-success {
        background: #dff0d8;
        color: #333;
      }
      h1 small {
        font-size:65%;
        color:#777;
      }
      nav {
        margin-bottom: 1em;
      }
    </style>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/awesomplete/1.1.2/awesomplete.min.css" />
    {{ template "style" . }}

    <!--[if lt IE 9]>
      <script src="http://html5shiv.googlecode.com/svn/trunk/html5.js"></script>
    <![endif]-->
  </head>

  <body>
    <h1>Tagaa <small>{{printv .Version}}</small></h1>
    {{ template "content" . }}
    <script src="https://cdnjs.cloudflare.com/ajax/libs/awesomplete/1.1.2/awesomplete.min.js"></script>
    {{ template "script" . }}
  </body>
</html>
{{ end }}
{{ define "style" }}{{end}}
{{ define "script" }}{{end}}
`

	indexTemplate = `
{{ define "style" }}
  <style>
    #advanced {
      display: none;
    }
    .image {
      max-width: 100%;
    }
  </style>
{{ end }}
{{ define "content" }}
  <nav>
    <a href="/upload">Upload</a>
  </nav>

  {{ $inputSize := 60 }}
  {{ $taRows := 6 }}

  {{ if .Err }}
    <div class="block block-danger">
      {{ .Err }}
    </div>
  {{ end }}
  <form action="/load" method="POST" enctype="multipart/form-data">
    <label for="loadCSVFile"><b>Load CSV File</b></label>
    <br>
    <input id="loadCSVFile" name="csvFilename" type="file" accept=".csv" required>
    <input type="submit" value="Load from CSV">
    <button id="toggleButton" type="button" onclick="toggleAdvanced()">Advanced +</button>
    <br>
  </form>
  <form action="/update" method="POST">
    <div id="advanced">
      <label for="csvFilenameInput"><b>CSV Filename</b></label>
      <br>
      <input id="csvFilenameInput" type="text" name="csvFilename" value="{{ .CSVFilename }}" size="{{ $inputSize }}">
      <input id="saveCSVSubmit" type="submit" value="Save to CSV">
      <br>
      <label for="directory"><b>Working Directory</b></label>
      <br>
      <input id="directory" type="text" name="prefix" value="{{ .WorkingDir }}" disabled size="{{ $inputSize }}">
      <br>
      <label for="prefixInput"><b>Server Path Prefix</b> (It will replace working directory path prefix)</label>
      <br>
      <input id="prefixInput" type="text" name="prefix" value="{{ .Prefix }}" size="{{ $inputSize }}">
      <br>
      <label for="useLinuxSepInput"><b>Use Linux Separator "/" when saving to CSV</b> </label>
      <br>
      <input id="useLinuxSepInput" type="checkbox" name="useLinuxSep" {{if eq .UseLinuxSep true}}checked{{end}}>
      (Check, if working on a windows machine and want to upload to a Linux machine)
      <input id="scroll" type="hidden" name="scroll" value="">
    </div>

    <section>
      {{ if .Images }}
        <h2>Images</h2>
      {{ else }}
        <h2>No Images found in local directory</h2>
        Add some and then refresh.
      {{ end }}

      {{ range .Images }}
        <article>
          <fieldset>
            <a id="tags{{ .ID }}"></a>
            <legend>{{ .Name }}</legend>
            <a href="#img{{ .ID }}"><img class="image" src="/img/{{ .ID }}" alt="{{ .Name }}"></a>
            <br>
            <label for="tagsTextArea{{ .ID }}"><b>Tags</b></label>
            <br>
            <textarea id="tagsTextArea{{ .ID }}" name="image[{{ .ID }}].tags" class="awesomeplete" data-multiple cols="{{ $inputSize }}" rows="{{ $taRows }}">{{ join .Tags " " }}</textarea>
            <br>
            <label for="sourceInput{{ .ID }}"><b>Source</b></label>
            <br>
            <input id="sourceInput{{ .ID }}" type="text" name="image[{{ .ID }}].source" value="{{ .Source }}" size="{{ $inputSize }}">
            <br>
            <label><b>Rating</b></label>
            <br>
            <input id="sRadio{{ .ID }}" type="radio" name="image[{{ .ID }}].rating" value="s" {{ if eq .Rating "s" }}checked{{ end }}>
            <label for="sRadio{{ .ID }}">Safe</label>
            <input id="qRadio{{ .ID }}" type="radio" name="image[{{ .ID }}].rating" value="q" {{ if eq .Rating "q" }}checked{{ end }}>
            <label for="qRadio{{ .ID }}">Questionable</label>
            <input id="eRadio{{ .ID }}" type="radio" name="image[{{ .ID }}].rating" value="e" {{ if eq .Rating "e" }}checked{{ end }}>
            <label for="eRadio{{ .ID }}">Explicit</label>
            <br>
            <input type="submit" value="Save to CSV" onclick="setScroll(this)" data-scroll="#tags{{.ID}}">
            <a id="img{{ .ID }}"></a>
          </fieldset>
        </article>
        <br>
      {{ end }}
    </section>
  </form>
{{ end }}
{{ define "script" }}
  <script>
    (function(){
      "use strict";

      function setScroll(e) {
        var scroll = e.getAttribute("data-scroll");
        document.getElementById("scroll").value = scroll;
      }

      function toggleAdvanced() {
        b = document.getElementById("toggleButton");
        div = document.getElementById("advanced");
        // Empty display reverts to CSS rule, in this case none.
        if (div.style.display == '') {
          div.style.display = 'block';
          b.innerHTML = "Advanced -";
        } else {
          div.style.display = '';
          b.innerHTML = "Advanced +";
        }
      }

      // Autocomplete

      var map = {};
      var tas = document.querySelectorAll('textarea[data-multiple]');
      tas.forEach(function(ta){
        var ap = makeAwesomplete(ta);
        map[ta.id] = ap;
        ta.onkeyup = getTagsEventHandler;
      });
      function makeAwesomplete(ta) {
        return new Awesomplete(ta, {
          filter: function(text, input) {
            return Awesomplete.FILTER_CONTAINS(text, input.match(/[^ ]*$/)[0]);
          },

          item: function(text, input) {
            return Awesomplete.ITEM(text, input.match(/[^ ]*$/)[0]);
          },

          replace: function(text) {
            var before = this.input.value.match(/^.+ \s*|/)[0];
            this.input.value = before + text.value + " ";
          },
          // Set sort function to false to disable sorting. Our backend handler
          // returns items sorted by count (first kusubooru then danbooru).
          sort: false
        });
      }

      var timeout = null;
      function getTagsEventHandler(e) {
        var code = (e.keyCode || e.which);
        // https://github.com/LeaVerou/awesomplete/issues/16802#issuecomment-303124988
        if (code !== 37 && code !== 38 && code !== 39 && code !== 40 && code !== 27 && code !== 13) {
          var input = this.value;
          var id = this.id;
          // Wait for user to stop typing before getting tags:
          // https://schier.co/blog/2014/12/08/wait-for-user-to-stop-typing-using-javascript.html
          clearTimeout(timeout);

          timeout = setTimeout(function () {
              getTags(input.match(/[^ ]*$/)[0], id);
          }, 500);
        }
      }

      function getTags(query, apid) {
        if (query == "" || query.length < 3) {
          return;
        }
        var list=[];
        var xhr = new XMLHttpRequest();
        xhr.onreadystatechange = function(response) {
          if (xhr.readyState === 4) {
            if (xhr.status === 200) {
              var tags = JSON.parse(xhr.responseText);
              tags.forEach(function(item) {
                var label = item.name;
                if (item.old) {
                  label = item.old+" → "+item.name;
                }
                if (item.category == "kusubooru") {
                  label = '<img src="img/kusubooru.ico" style="float:left;margin-right:2px;height:16px;width:16px">' + label
                }
                if (item.category == "danbooru") {
                  label = '<img src="img/danbooru.ico" style="float:left;margin-right:2px">' + label
                }
                label = label + '<span style="float:right">'+item.count+'</span>';
                list.push({"label": label, "value": item.name, "cound": item.count});
              });
              map[apid].list = list;
              // Update the placeholder text.
              //input.placeholder = "e.g. datalist";
            } else {
              // An error occured :(
              //input.placeholder = "Couldn't load datalist options :(";
            }
          }
        };
        xhr.open("GET", "tags?q="+query, true);
        xhr.send();
      }

    })();
  </script>
{{ end }}
`
	uploadTemplate = `
{{ define "style" }}
  <style>
    .thumbnail {
      width: 50px;
      height: 30px;
      display: inline-block;
    }
    .thumbnail img {
      width: 100%;
      height: auto;
    }
    .upload-table {
      width: 100%;
    }
    .upload-table textarea {
      width: 95%;
    }
    .upload-button {
      display: inline-block;
      padding: 0.5em;
    }

    .loader {
      display: none;
      border: 5px solid #f3f3f3;
      border-radius: 50%;
      border-top: 5px solid #006FFA;
      border-right: 5px solid #006FFA;
      width: 32px;
      height: 32px;
      -webkit-animation: spin 1s linear infinite;
      animation: spin 1s linear infinite;
      will-change: transform;
    }
    .loader-small {
      width: 8px;
      height: 8px;
      border-width: 3px;
    }
    @-webkit-keyframes spin {
      0% { -webkit-transform: rotate(0deg); }
      100% { -webkit-transform: rotate(360deg); }
    }
    @keyframes spin {
      0% { transform: rotate(0deg); }
      100% { transform: rotate(360deg); }
    }
  </style>
{{end}}

{{ define "content" }}
  <nav>
    <a href="/">Back</a>
  </nav>

  {{ if .Err }}
    <div class="block block-danger">
      {{ .Err }}
    </div>
  {{ else if .Success }}
    <div class="block block-success">
     {{ .Success }}
    </div>
  {{ end }}

  <form action="/upload" method="POST" enctype="multipart/form-data" onsubmit="showLoader()">
    <table class="upload-table">
      <thead>
        <tr>
          <th></th>
          <th>Name</th>
          <th>Tags</th>
          <th>Source</th>
          <th>Rating</th>
        </tr>
      </thead>
      <tbody>
        {{ range .Images }}
          <tr>
            <td>
              <div class="thumbnail">
                <a href="#img{{ .ID }}"><img src="/img/{{ .ID }}" alt="{{ .Name }}" width=150 height=100></a>
              </div>
            </td>
            <td width="10%">
              {{ .Name }}
            </td>
            <td width="65%">
              <textarea id="tagsTextArea{{ .ID }}" name="image[{{ .ID }}].tags" cols="20" rows="2" readonly>{{ join .Tags " " }}</textarea>
            </td>
            <td width="25%">
              {{ .Source }}
            </td>
            <td>
              {{ if eq .Rating "s" }} Safe
              {{ else if eq .Rating "q" }} Questionable
              {{ else if eq .Rating "e" }} Explicit
              {{ else }} Unknown
              {{ end }}
            </td>
          </tr>
        {{ end }}
      </tbody>
    </table>

    <span>The images above are going to be:</span>
    <ul>
      <li>Compressed to a .zip archive</li>
      <li>Uploaded to server Kusubooru.com</li>
      <li>Manually reviewed before posted</li>
    </ul>
    <p>Please make sure that all images have adequate tags, a source and a rating before uploading.</p>

    <p>Use your Kusubooru account to upload:</p>
    <label for="username">Username</label>
    <input id="username" type="text" name="username" placeholder="Username" required>
    <label for="password">Password</label>
    <input id="password" type="password" name="password" placeholder="Password" required>
    <button type="button" onclick="testCredentials()">Test</button>
    <label id="result"></label>
    <div id="testLoader" class="loader loader-small"></div>

    <p><small>(Max file size for a single upload is 50MB and you may upload a total of 200MB per day.)</small></p>
    <input id="uploadButton" class="upload-button" type="submit" value="Upload">
    <div id="uploadLoader" class="loader"></div>
  </form>
{{ end }}
{{ define "script" }}
  <script>
    function showLoader() {
      var loader = document.getElementById("uploadLoader");
      loader.style.display = "inline-block";
      var button = document.getElementById("uploadButton");
      button.style.display = "none";
    }

    function testCredentials() {
      var loader = document.getElementById("testLoader");
      loader.style.display = "inline-block";
      var username = document.getElementById("username").value;
      var password = document.getElementById("password").value;
      var resultLabel = document.getElementById("result");
      resultLabel.innerHTML = "";
      var xhr = new XMLHttpRequest();
      var url = "https://kusubooru.com/suggest/login/test";
      var params = "username="+username+"&password="+password;
      xhr.open("POST", url, true);
      xhr.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
      xhr.onreadystatechange = function() {
        loader.style.display = "none";
        if(xhr.readyState == 4 && xhr.status == 200) {
          resultLabel.innerHTML = "Ok!";
	      } else if(xhr.readyState == 4 && xhr.status != 200) {
	        var reason = "";
	        if (xhr.responseText) { reason = ": " + xhr.responseText }
          resultLabel.innerHTML = "Failed" + reason;
        }
      }
      xhr.send(params);
    };
  </script>
{{ end }}
`
)

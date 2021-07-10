package main

const index = `<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, minimal-ui">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/9.18.1/styles/default.min.css">
    <script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/9.18.1/highlight.min.js"></script>
    <script>hljs.initHighlightingOnLoad();</script>
    <style>
      {{ .Stylesheet }}
    </style>
    <style>
     body {
        box-sizing: border-box;
        min-width: 200px;
        max-width: 980px;
        margin: 0 auto;
        padding: 45px;
      }
     @media (max-width: 767px) {
       .markdown-body {
         padding: 15px;
       }
     }
	  </style>
  </head>
  <body>
    <article class="markdown-body">
      {{ .Content }}
    </article>
  </body>
</html>
`

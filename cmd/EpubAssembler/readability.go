// This code was automatically generated since I do not know Javascript
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dop251/goja"
)

func NewJSWorkerFactory() (*JSWorkerFactory, error) {
	program, err := goja.Compile("readability.js", readabilityJS, false)
	if err != nil {
		return nil, fmt.Errorf("failed to compile readability.js: %w", err)
	}
	return &JSWorkerFactory{readabilityProgram: program}, nil
}

func (f *JSWorkerFactory) NewJSWorker() (*JSWorker, error) {
	vm := goja.New()
	err := vm.Set("window", vm.NewObject())
	if err != nil {
		return nil, fmt.Errorf("failed to set window: %w", err)
	}

	// Create a basic document object
	_, err = vm.RunString(`
		window.document = {
			createElement: function(tag) {
				return {
					tagName: tag.toUpperCase(),
					appendChild: function(child) {},
					setAttribute: function(name, value) {},
					getAttribute: function(name) { return null; },
					innerHTML: '',
					textContent: '',
					childNodes: [],
					parentNode: null
				};
			},
			createTextNode: function(text) {
				return { textContent: text, nodeType: 3 };
			},
			implementation: {
				createHTMLDocument: function() {
					return window.document;
				}
			}
		};
		
		// Basic console for debugging
		var console = {
			log: function() {},
			warn: function() {},
			error: function() {}
		};
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to setup DOM environment: %w", err)
	}

	readabilityJS = loadReadabilityJS()
	if err != nil {
		return nil, fmt.Errorf("failed to load readability.js: %w", err)
	}

	_, err = vm.RunProgram(f.readabilityProgram)
	if err != nil {
		return nil, fmt.Errorf("failed to execute readability.js: %w", err)
	}

	return &JSWorker{vm: vm}, nil
}

func (worker *JSWorker) ParseHTML(htmlContent HTML) (*ReadabilityObject, error) {
	parseScript := `
		function parseWithReadability(htmlString) {
			// Create a basic DOM parser
			var doc = {
				documentElement: null,
				createElement: function(tag) {
					return {
						tagName: tag.toUpperCase(),
						appendChild: function(child) {},
						setAttribute: function(name, value) {},
						getAttribute: function(name) { return null; },
						innerHTML: '',
						textContent: '',
						childNodes: [],
						parentNode: null,
						querySelectorAll: function() { return []; },
						querySelector: function() { return null; }
					};
				},
				querySelectorAll: function() { return []; },
				querySelector: function() { return null; }
			};
			
			// Very basic HTML parsing - in production you'd want something more robust
			// This creates a minimal document structure
			var titleMatch = htmlString.match(/<title[^>]*>([^<]*)<\/title>/i);
			var title = titleMatch ? titleMatch[1] : '';
			
			// Create article element with the HTML content
			var article = {
				innerHTML: htmlString,
				textContent: htmlString.replace(/<[^>]+>/g, ' ').replace(/\s+/g, ' ').trim(),
				getAttribute: function(name) {
					if (name === 'class') return '';
					return null;
				},
				querySelectorAll: function() { return []; },
				querySelector: function() { return null; },
				childNodes: [],
				parentNode: doc
			};
			
			doc.documentElement = article;
			
			try {
				var reader = new Readability(doc, {
					debug: false,
					maxElemsToParse: 0,
					nbTopCandidates: 5,
					charThreshold: 500,
					classesToPreserve: []
				});
				
				var result = reader.parse();
				
				if (result) {
					return {
						title: result.title || title,
						content: result.content || '',
						excerpt: result.excerpt || ''
					};
				}
				
				// Fallback if readability fails
				return {
					title: title,
					content: htmlString,
					excerpt: ''
				};
				
			} catch (e) {
				// Fallback parsing
				var bodyMatch = htmlString.match(/<body[^>]*>([\s\S]*)<\/body>/i);
				var bodyContent = bodyMatch ? bodyMatch[1] : htmlString;
				var textContent = bodyContent.replace(/<[^>]+>/g, ' ').replace(/\s+/g, ' ').trim();
				
				return {
					title: title,
					content: textContent,
					excerpt: textContent.substring(0, 200)
				};
			}
		}
	`

	_, err := worker.vm.RunString(parseScript)
	if err != nil {
		return nil, fmt.Errorf("failed to create parser function: %w", err)
	}

	err = worker.vm.Set("htmlInput", htmlContent)
	if err != nil {
		return nil, fmt.Errorf("failed to set HTML input: %w", err)
	}

	val, err := worker.vm.RunString("parseWithReadability(htmlInput)")
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	resultJSON, err := json.Marshal(val.Export())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %w", err)
	}

	var result ReadabilityObject
	err = json.Unmarshal(resultJSON, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal result: %w", err)
	}

	return &result, nil
}

func (worker *JSWorker) ParseURL(url URL) (*ReadabilityObject, error) {
	resp, err := http.Get(string(url))
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	result, err := worker.ParseHTML(HTML(body))
	if err != nil {
		return nil, err
	}
	return result, nil
}

func loadReadabilityJS() string {
	return readabilityJS
}

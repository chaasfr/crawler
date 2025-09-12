package main

import (
	"net/url"
	"reflect"
	"testing"
)

func TestGetH1FromHTML(t *testing.T) {
	tests := []struct {
		name          string
		inputBody     string
		expected      []string
	}{
		{
			name:     "one H1",
			inputBody: `
			<html>
				<body>
			  		<h1>Welcome to Boot.dev</h1>
					<main>
						<p>Learn to code by building real projects.</p>
						<p>This is the second paragraph.</p>
					</main>
				</body>
		  	</html>
			`,
			expected: []string{"Welcome to Boot.dev"},
		},
		{
			name:     "empty results",
			inputBody: `
		<html>
			<body>
			</body>
		</html>
		`,
			expected: []string{},
		},
		{
			name:     "two H1 no H2",
			inputBody: `
			<html>
				<body>
			  		<h1>Welcome to Boot.dev</h1>
					<main>
						<p>Learn to code by building real projects.</p>
						<p>This is the second paragraph.</p>
					</main>
					<h1>Hope you like it!</h1>
					<h2>Truly, I do</h2>
				</body>
		  	</html>
			`,
			expected: []string{"Welcome to Boot.dev", "Hope you like it!"},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := getH1FromHTML(tc.inputBody)
			if err != nil {
				t.Fatalf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if !reflect.DeepEqual(tc.expected, actual) {
				t.Errorf("Test %v - %s FAIL: expected Headers: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestGetFirstParagraphFromHTMLMainPriority(t *testing.T) {
	tests := []struct {
		name          string
		inputBody     string
		expected      string
	}{
		{
			name:     "one paragraph",
			inputBody: `
			<html>
				<body>
			  		<h1>Welcome to Boot.dev</h1>
					<main>
						<p>Learn to code by building real projects.</p>
					</main>
				</body>
		  	</html>
			`,
			expected: "Learn to code by building real projects.",
		},
		{
			name:     "empty results",
			inputBody: `
		<html>
			<body>
			</body>
		</html>
		`,
			expected: "",
		},
		{
			name:     "two p one main",
			inputBody: `
			<html>
				<body>
					<p>This is the second paragraph.</p>
			  		<h1>Welcome to Boot.dev</h1>
					<main>
						<p>Learn to code by building real projects.</p>
					</main>
					<h1>Hope you like it!</h1>
					<h2>Truly, I do</h2>
				</body>
		  	</html>
			`,
			expected: "Learn to code by building real projects.",
		},
		{
			name:     "one p no main",
			inputBody: `
			<html>
				<body>
			  		<h1>Welcome to Boot.dev</h1>
					<p>Learn to code by building real projects.</p>
					<h1>Hope you like it!</h1>
					<h2>Truly, I do</h2>
				</body>
		  	</html>
			`,
			expected: "Learn to code by building real projects.",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := getFirstParagraphFromHTML(tc.inputBody)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected 1st Paragraph: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestGetURLFromHTML(t *testing.T) {
	tests := []struct {
		name          string
		inputURL      string
		inputBody     string
		expected      []string
	}{
		{
			name:     "absolute and relative URLs",
			inputURL: "https://blog.boot.dev",
			inputBody: `
		<html>
			<body>
				<a href="/path/one">
					<span>Boot.dev</span>
				</a>
				<a href="https://other.com/path/one">
					<span>Boot.dev</span>
				</a>
			</body>
		</html>
		`,
			expected: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
		},
		{
			name:     "empty results",
			inputURL: "https://blog.boot.dev",
			inputBody: `
		<html>
			<body>
			</body>
		</html>
		`,
			expected: []string{},
		},
		{
			name:     "multiple bodies URLs",
			inputURL: "https://blog.boot.dev",
			inputBody: `
		<html>
			<body>
				<a href="/path/one">
					<span>Boot.dev</span>
				</a>
			</body>
			<body>
				<a href="/path/two/one">
					<span>Boot.dev</span>
				</a>
			</body>
		</html>
		`,
			expected: []string{"https://blog.boot.dev/path/one", "https://blog.boot.dev/path/two/one"},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			urlNormalized, err := url.Parse(tc.inputURL)
			if err != nil {
				t.Errorf("couldn't parse input URL: %v", err)
				return
			}
			actual, err := getURLsFromHTML(tc.inputBody, urlNormalized)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if !reflect.DeepEqual(tc.expected, actual) {
				t.Errorf("Test %v - %s FAIL: expected URLs: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestGetImageFromHTML(t *testing.T) {
	tests := []struct {
		name          string
		inputURL      string
		inputBody     string
		expected      []string
	}{
		{
			name:     "relative URLs",
			inputURL: "https://blog.boot.dev",
			inputBody: `
		<html>
			<body>
				<img src="/logo.png" alt="Logo">
			</body>
		</html>
		`,
			expected: []string{"https://blog.boot.dev/logo.png"},
		},{
			name:     "absolute URLs",
			inputURL: "https://blog.boot.dev",
			inputBody: `
		<html>
			<body>
				<img src="https;//www.mywebsite.com/logo.png" alt="Logo">
			</body>
		</html>
		`,
			expected: []string{"https;//www.mywebsite.com/logo.png"},
		},{
			name:     "no image",
			inputURL: "https://blog.boot.dev",
			inputBody: `
		<html>
			<body>
			</body>
		</html>
		`,
			expected: []string{},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			urlNormalized, err := url.Parse(tc.inputURL)
			if err != nil {
				t.Errorf("couldn't parse input URL: %v", err)
				return
			}
			actual, err := getImagesFromHTML(tc.inputBody, urlNormalized)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if !reflect.DeepEqual(tc.expected, actual) {
				t.Errorf("Test %v - %s FAIL: expected URLs: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
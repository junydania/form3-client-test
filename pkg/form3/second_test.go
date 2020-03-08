
const (
	okResponse = `{
		"users": [
			{"id": 1, "name": "Roman"},
			{"id": 2, "name": "Dmitry"}
		]	
	}`
)

func TestServe(t *testing.T) {
    // The method to use if you want to practice typing
    s := &http.Server{
        Handler: http.HandlerFunc(ServeHTTP),
    }
    // Pick port automatically for parallel tests and to avoid conflicts
    l, err := net.Listen("tcp", ":0")
    if err != nil {
        t.Fatal(err)
    }
    defer l.Close()
    go s.Serve(l)

    res, err := http.Get("http://" + l.Addr().String() + "/?sloths=arecool")
    if err != nil {
        log.Fatal(err)
    }
    greeting, err := ioutil.ReadAll(res.Body)
    res.Body.Close()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(greeting))
}


func ItemTest(t *testing.T) {
	// Common setup.  Abstract this out
	// This allows each test to create its own handler by changing handler variable
	handler := http.NotFound
	hs := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		handler(rw, req)
	}))
	defer hs.Close()
	// Notice I set the base URL of the client to the httptest server
	c := Client{
		BaseURL: hs.URL,
	}

	// Code specific to this test
	handler = func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/v1/smite/item/sword" {
			t.Error("Bad path!")
		}
		io.WriteString(rw, `{"type":"sword"}`)
	}
	item, err := c.GetItem(ctx, "sword")
	if err != nil {
		t.Error("Got error sending item")
	}
	if item.Type != "sword" {
		t.Error("Did not get a sword!")
	}
}

func TestServer(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Send response to be tested
		rw.Write([]byte(`OK`))
	}))
	// Close the server when test finishes
	defer server.Close()

	//Use Client & URL from our local test server
	api := Client{server.Client(), server.URL}
	body, err := api.Start()

	ok(t, err)
	equals(t, []byte("OK"), body)

}


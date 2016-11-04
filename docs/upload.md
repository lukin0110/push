# Upload snippets

* [POST File](http://stackoverflow.com/questions/20205796/golang-post-data-using-the-content-type-multipart-form-data)
* [Upload with Forms](https://github.com/astaxie/build-web-application-with-golang/blob/master/en/04.5.md)
* [Upload Binary](http://stackoverflow.com/questions/29529926/http-post-data-binary-curl-equivalent-in-golang)

## Snippet

```go
func Upload(url, file string) (string, error) {
    // Prepare a form that you will submit to that URL.
    var b bytes.Buffer
    w := multipart.NewWriter(&b)
    // Add your image file
    f, err := os.Open(file)

    if err != nil {
        return "", err
    }
    defer f.Close()
    fw, err := w.CreateFormFile("filename", file)

    if err != nil {
        return "", err
    }

    if _, err = io.Copy(fw, f); err != nil {
        return "", err
    }

    // Don't forget to close the multipart writer.
    // If you don't close it, your request will be missing the terminating boundary.
    w.Close()

    // Now that you have a form, you can submit it to your handler.
    req, err := http.NewRequest("PUT", url, &b)
    if err != nil {
        return "", err
    }

    // Don't forget to set the content type, this will contain the boundary.
    //req.Header.Set("Content-Type", w.FormDataContentType())
    //text/markdown; charset=UTF-8
    req.Header.Set("Content-Type", "text/markdown; charset=UTF-8")

    // Submit the request
    client := &http.Client{}
    res, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer res.Body.Close()

    // Check the response
    if res.StatusCode != http.StatusOK {
        err = fmt.Errorf("bad status: %s", res.Status)
    } else {
	bs, err1 := ioutil.ReadAll(res.Body)
	if(err1 != nil) {
	    return "", err
	}
	return string(bs), nil
    }

    return "", nil
}
```
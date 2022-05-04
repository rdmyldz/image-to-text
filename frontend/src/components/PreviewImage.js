import React, { useEffect, useRef, useState } from "react";
import { Button, ButtonGroup, Col, Image, Spinner } from "react-bootstrap";

const PreviewImage = () => {
  const [image, setImage] = useState(null);
  const [preview, setPreview] = useState();
  const [data, setData] = useState(null);
  const [copyButton, setCopyButton] = useState("Copy");
  const [spinnerVisible, setSpinnerVisible] = useState(false);
  const fileRef = useRef();

  useEffect(() => {
    if (image) {
      const reader = new FileReader();
      reader.onloadend = () => {
        setPreview(reader.result);
      };
      reader.readAsDataURL(image);
    } else {
      setPreview("");
    }
  }, [image]);

  const handleConvert = (e) => {
    setSpinnerVisible(true);
    e.preventDefault();
    const url = "/upload";
    const formData = new FormData();
    formData.append("file", image);
    const requestOptions = {
      method: "POST",
      body: formData,
    };
    fetch(url, requestOptions)
      .then((res) => res.json())
      .then((result) => {
        setData(result);
        setSpinnerVisible(false);
      })
      .catch((error) => console.log("Form submit error", error));
  };

  const handleCopy = async (text) => {
    await navigator.clipboard.writeText(text);
    setCopyButton("Copied");
    setTimeout(() => {
      setCopyButton("Copy");
    }, 1000);
  };
  return (
    <>
      <Col className="col-12 text-center ">
        <form>
          <ButtonGroup className="gap-3 ">
            <Button
              role="group"
              type="button"
              onClick={() => {
                fileRef.current.click();
              }}
            >
              Add Image
            </Button>
            <Button role="group" onClick={handleConvert}>
              Convert Image
            </Button>
          </ButtonGroup>
          <input
            type="file"
            name="file"
            style={{ display: "none" }}
            ref={fileRef}
            accept="image/*"
            onChange={(e) => {
              const file = e.target.files[0];
              if (file && file.type.includes("image")) {
                setImage(file);
              } else {
                setImage(null);
              }
              setData(null);
            }}
          />
        </form>
      </Col>
      {spinnerVisible && !data && (
        <Col className="col-12 text-center">
          <Spinner animation="border" variant="primary" role="status">
            <span className="visually-hidden">Loading...</span>
          </Spinner>
        </Col>
      )}
      <Col className="col-12 col-md-6 p-3 align-self-stretch mw-50">
        {image && (
          <Image
            fluid
            className="d-block ms-auto me-auto"
            src={preview}
            alt={image.name}
          />
        )}
      </Col>
      <Col className="col-12 col-md-6 position-relative p-3 pt-3 align-self-stretch mw-50">
        {data && (
          <>
            <div>{data.content}</div>
            <Button
              className="copy-button btn-clipboard"
              data-bs-toggle="tooltip"
              data-bs-placement="top"
              title="Copy to clipboard"
              onClick={() => handleCopy(data.content)}
            >
              {copyButton}
            </Button>
          </>
        )}
      </Col>
    </>
  );
};

/*

Breakpoint			Class-infix		Dimensions
X-Small					None			<576px
Small					sm				≥576px
Medium					md				≥768px
Large					lg				≥992px
Extra 					large			xl	≥1200px
Extra extra large		xxl				≥1400px
*/

export default PreviewImage;

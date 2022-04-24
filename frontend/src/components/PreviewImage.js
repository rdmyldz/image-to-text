import React, { useEffect, useRef, useState } from "react";
import { Button, Col, Image, Row } from "react-bootstrap";

const PreviewImage = () => {
  const [image, setImage] = useState(null);
  const [preview, setPreview] = useState();
  const [data, setData] = useState(null);
  const [copyButton, setCopyButton] = useState("Copy");
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
      .then((result) => setData(result))
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
      <Row className="py-3">
        <Col className="justify-content-center d-flex">
          <form>
            <button
              type="button"
              onClick={(e) => {
                // e.preventDefault();
                fileRef.current.click();
              }}
            >
              Add Image
            </button>
            <button onClick={handleConvert}>Convert Image</button>
            <input
              type="file"
              name="file"
              style={{ display: "none" }}
              ref={fileRef}
              //   accept="image/*"
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
      </Row>
      <Row className="row-cols-1 row-cols-lg-2 p-3  ">
        {/* <div className="  image-div"> */}
        {image && (
          <Col className="col-lg-6 border border-primary p-3 ">
            {" "}
            <Image fluid src={preview} alt={image.name} />
          </Col>
        )}
        {/* </div> */}
        {data && (
          <Col className="col-lg-6 border border-primary p-4 position-relative">
            <div className="mt-3">{data.content}</div>
            <Button
              className="copy-button btn-clipboard"
              data-bs-toggle="tooltip"
              data-bs-placement="top"
              title="Copy to clipboard"
              onClick={() => handleCopy(data.content)}
            >
              {copyButton}
            </Button>
          </Col>
        )}
      </Row>
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

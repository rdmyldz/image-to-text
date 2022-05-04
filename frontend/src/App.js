import { Container, Row } from "react-bootstrap";
import PreviewImage from "./components/PreviewImage";
import Title from "./components/Title";
import "bootstrap/dist/css/bootstrap.min.css";
import "./App.css";

// TODO: make Convert Button disabled when pushed once, after fixing backend consecutive request issue

function App() {
  return (
    <Container fluid>
      <Row className="justify-content-center align-items-center gy-3">
        <Title />
        {/* <UploadForm /> */}
        <PreviewImage />
      </Row>
    </Container>
  );
}

export default App;

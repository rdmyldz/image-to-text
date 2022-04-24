import { Container, Row } from "react-bootstrap";
import PreviewImage from "./components/PreviewImage";
import Title from "./components/Title";
import UploadForm from "./components/UploadForm";
import "bootstrap/dist/css/bootstrap.min.css";
import "./App.css";

function App() {
  return (
    <Container fluid className="mt-3">
      <Title className="d-flex justify-content-center" />
      {/* <UploadForm /> */}
      <PreviewImage />
    </Container>
  );
}

export default App;

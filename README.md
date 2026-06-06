# Word Templates API (Golang)

A high-performance REST API built in Golang to generate PDF documents dynamically by injecting JSON data into Microsoft Word (`.docx`) templates and converting them via LibreOffice.

## Architecture

*   **Language:** Golang (using the Gin framework)
*   **Templating:** Native ZIP/XML parser (handles spaces inside curly braces, e.g., `{{ key }}`)
*   **PDF Conversion:** Headless LibreOffice
*   **Deployment:** Multi-stage Docker container (Alpine-based)

---

## How to Run

### Option 1: Running with Docker Compose (Recommended)

Docker Compose is the easiest way to run the application, as it packages Go and LibreOffice together. The `compose.yml` mounts the `docs/` and `fonts/` directories, so you can add or update templates/fonts without rebuilding the container.

1. Start the service:
   ```bash
   docker compose up --build
   ```
2. The server will start and listen on `http://localhost:5000`.

### Option 2: Running Locally (Requires Go & LibreOffice)

To run the application locally without Docker, you must have **Go** (1.18+) and **LibreOffice** installed on your system.

1. Ensure `libreoffice` is available in your system path (e.g., `libreoffice --version`).
2. Run the application:
   ```bash
   go run .
   ```
3. The server will start on port `5000` (or the port defined by the `PORT` environment variable).

---

## API Documentation & Postman Integration

### Endpoint: `POST /generate`
*   **URL:** `http://localhost:5000/generate`
*   **Headers:** `Content-Type: application/json`
*   **Description:** Takes a flat JSON payload. The `template` parameter specifies the target `.docx` file in the `docs/` folder. Any other key-value pairs are treated as variables to be replaced inside the template.

---

### Importing to Postman

We have provided a pre-configured Postman Collection that contains all the endpoint requests and examples. You can import it directly to test everything.

#### Step-by-Step Instructions:

1. **Import the Collection:**
   - Open **Postman**.
   - Click the **Import** button in the top-left corner (or use `Ctrl + O` / `Cmd + O`).
   - Select **File** and choose the `word-templates.postman_collection.json` file from the root directory of this project.
   - Click **Import** to add the "Word Templates API" collection to your workspace.

2. **Run the Requests:**
   - In the left sidebar, expand the **Word Templates API** collection.
   - You will see three pre-configured requests:
     *   `Generate PDF - Carta de Invitación`
     *   `Generate PDF - Certificado Docente`
     *   `Generate PDF - Cronograma de Clases`
   - Select any request to view its body and settings.

3. **Download the PDF Response:**
   - Since the API returns a binary PDF file, clicking the standard **Send** button will show binary text in the response pane.
   - To view the generated PDF, click the **arrow next to the "Send" button** and select **"Send and Download"**.
   - Save the file with a `.pdf` extension (e.g., `output.pdf`) and open it using any PDF viewer.

---

## Example JSON Payloads

If you prefer to configure the requests manually, use the following payloads:

### 1. Carta de Invitación
*   **Template:** `carta-invitacion.docx`
*   **Body:**
    ```json
    {
      "template": "carta-invitacion.docx",
      "hoy": "26 de mayo de 2026",
      "docente": "Juan Perez",
      "nombre_diplomado": "DIPLOMADO EN EDUCACIÓN SUPERIOR",
      "nombre_modulo": "METODOLOGÍAS ÁGILES",
      "competencia_modulo": "Aplica metodologías ágiles para el desarrollo de proyectos.",
      "contenidos_minimos": "1. Introducción\n2. Scrum\n3. Kanban\n4. Evaluación",
      "dias_clases": "Lunes y Miércoles",
      "fecha_clase_1": "01/06/2026",
      "fecha_clase_2": "03/06/2026",
      "fecha_clase_3": "08/06/2026",
      "fecha_clase_4": "10/06/2026",
      "fecha_clase_5": "15/06/2026",
      "fecha_clase_6": "17/06/2026",
      "objetivo": "Mejorar las competencias docentes a través de nuevas metodologías."
    }
    ```

### 2. Certificado Docente
*   **Template:** `certificado-docente.docx`
*   **Body:**
    ```json
    {
      "template": "certificado-docente.docx",
      "hoy": "30 de junio de 2026",
      "docente": "María Antonieta",
      "fecha_inicio": "01/05/2026",
      "fecha_fin": "30/06/2026"
    }
    ```

### 3. Cronograma de Clases
*   **Template:** `cronograma-clases.docx`
*   **Body:**
    ```json
    {
      "template": "cronograma-clases.docx",
      "nombre_modulo": "Desarrollo Backend Avanzado",
      "dias_clases": "Martes y Jueves",
      "docente": "Ing. Carlos Ruiz",
      "fecha_clase_1": "10/08/2026",
      "fecha_clase_2": "12/08/2026",
      "fecha_clase_3": "17/08/2026",
      "fecha_clase_4": "19/08/2026",
      "fecha_clase_5": "24/08/2026",
      "fecha_clase_6": "26/08/2026"
    }
    ```

---

## Template Customization

*   **Template Location:** Place your Word `.docx` templates inside the `docs/` folder.
*   **Placeholders:** Inside the Word document, wrap your variables in double curly braces, e.g., `{{docente}}` or `{{ docente }}`.
*   **Newlines:** The API automatically replaces `\n` characters in your JSON values with native Word XML line breaks before rendering the PDF.
*   **Custom Fonts:** To use custom fonts, place them in the `fonts/` folder. They will be copied and registered in the system during Docker container build/run.

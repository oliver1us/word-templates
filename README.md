# Word Templates API (Golang)

A high-performance REST API built in Golang to generate PDF documents dynamically by injecting JSON data into Microsoft Word (`.docx`) templates.

## Architecture

*   **Language:** Golang (using Gin framework)
*   **Templating:** Native ZIP/XML high-performance parser (supports spacing inside braces `{{ key }}`)
*   **PDF Conversion:** LibreOffice headless mode
*   **Deployment:** Multi-stage Docker container (Alpine based)

## How to Run

The easiest way to run the service is using Docker Compose. The `compose.yml` mounts the `docs/` and `fonts/` directories so you can update templates without rebuilding the container.

```bash
docker compose up --build
```

The server will start on `http://localhost:5000`.

## API Documentation

### Generate PDF
*   **Endpoint:** `POST /generate`
*   **Headers:** `Content-Type: application/json`
*   **Description:** Takes a flat JSON payload. The `template` parameter specifies the target file name. Every other parameter in the JSON payload is dynamically treated as a template replacement variable.

#### Postman Example 1: Carta de Invitación

1.  Open Postman and create a new **POST** request to `http://localhost:5000/generate`
2.  Go to the **Headers** tab and add `Content-Type: application/json`
3.  Go to the **Body** tab, select **raw** (JSON format), and paste the following flat structure:

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

> **💡 Postman Tip:** Instead of clicking the standard "Send" button, click the small arrow next to it and select **"Send and Download"**. This will correctly download the binary response and save it as a `.pdf` file on your computer.

#### Postman Example 2: Certificado Docente

```json
{
  "template": "certificado-docente.docx",
  "hoy": "30 de junio de 2026",
  "docente": "María Antonieta",
  "fecha_inicio": "01/05/2026",
  "fecha_fin": "30/06/2026"
}
```

#### Postman Example 3: Cronograma de Clases

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

## Template Notes

*   Place your new `.docx` templates in the `docs/` folder.
*   Variables inside your Word document should be wrapped in double curly braces, e.g., `{{docente}}` or `{{ docente }}`.
*   The API automatically processes `\n` characters in the JSON strings and converts them to actual line breaks in the generated Word document before rendering the PDF.

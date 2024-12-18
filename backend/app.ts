import express from 'express';

const app = express();
const port = 3001; // Backend działa na innym porcie niż frontend

app.use(express.json());

app.get('/api/hello', (req, res) => {
  res.json({ message: 'Hello from the backend!' });
});

app.listen(port, () => {
  console.log(`Backend running at http://localhost:${port}`);
});
import express from "express";
import type { Express } from "express";

const app: Express = express();

app.use(express.json());

export default app;

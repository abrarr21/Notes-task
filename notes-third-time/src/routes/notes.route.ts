import { Router } from "express";
import type { IRouter } from "express";
import * as notesController from "../controllers/notes.controller.js";

const notesRouter: IRouter = Router();

// @route /api/notes
// @description Create a note with title and descritpion given by user
// @access Public
notesRouter.post("/", notesController.createNote);

// @route /api/notes
// @description Get all the notes
// @access Public
notesRouter.get("/", notesController.getAllNotes);

// @route /api/notes/:id
// @description Update the given note using id
// @access Public
notesRouter.patch("/:id", notesController.updateNote);

// @route /api/notes/:id
// @description Delete the given note using id
// @access Public
notesRouter.delete("/:id", notesController.deleteNote);

export default notesRouter;

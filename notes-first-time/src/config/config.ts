import dotenv from "dotenv";

dotenv.config();

if (!process.env.MONGODB_URI) {
  throw new Error("MONGODB_URI is not provided in .env file");
}

const config = {
  MONGODB_URI: process.env.MONGODB_URI,
};

export default config;

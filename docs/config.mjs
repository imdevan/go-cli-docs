const stage = process.env.NODE_ENV || "dev"
const isProduction = stage === "production"

export default {
  url: isProduction ? "https://devan.gg" : "http://localhost:4321",
  basePath: isProduction ? "/go-cli-docs" : "/",
  github: "https://github.com/imdevan/go-cli-docs/",
  githubDocs: "https://github.com/imdevan/go-cli-docs/",
  title: "go-cli-docs",
  description: "A go cli to build docs for your go cli.",
}

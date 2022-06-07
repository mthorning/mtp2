use clap::Parser;
use reqwest::blocking::{multipart::Form, Client};
use serde::Deserialize;
use std::path::PathBuf;
use std::io;

const API_URL: &str = "http://localhost:8000";

#[derive(Parser)]
struct Args {
    #[clap(subcommand)]
    action: Action,
}

#[derive(clap::Subcommand)]
enum Action {
    List,
    Add { path: PathBuf },
}

#[derive(Deserialize)]
struct Image {
    name: String,
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let args = Args::parse();
    match args.action {
        Action::List => {
            list_images()?;
        }
        Action::Add { path } => {
            add_image(path)?;
        }
    }

    Ok(())
}

fn list_images() -> Result<(), Box<dyn std::error::Error>> {
    let body: Vec<Image> = reqwest::blocking::get(format!("{}/images", API_URL))?.json()?;

    for image in body {
        println!("{}", image.name);
    }

    Ok(())
}

fn add_image(path: PathBuf) -> Result<(), Box<dyn std::error::Error>> {
    println!("What's it called?");
    let mut title = String::new();
    io::stdin().read_line(&mut title)?;

    let form = Form::new()
        .text("title", title.trim().to_string())
        .file("image", path)?;

    let client = Client::new();

    let res = client
        .post(format!("{}/add", API_URL))
        .multipart(form)
        .send()?;

    println!("success: {:?}", res);

    Ok(())
}

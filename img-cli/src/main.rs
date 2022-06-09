use clap::Parser;
use reqwest::blocking::{multipart::Form, Client};
use serde::Deserialize;
use std::io;
use std::{str, path::PathBuf};
use base64::{decode, encode};

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
    Remove,
}

#[derive(Deserialize, Ord, PartialOrd)]
struct Image {
    name: String,
}

impl PartialEq for Image {
    fn eq(&self, other: &Self) -> bool {
        self.name == other.name
    }
}
impl Eq for Image {}

fn main() {
    match run_program() {
        Err(e) => println!("Error: {}", e),
        _ => (),
    }
}

fn run_program() -> Result<(), Box<dyn std::error::Error>> {
    let args = Args::parse();
    match args.action {
        Action::List => {
            list_images()?;
        }
        Action::Add { path } => {
            add_image(&path)?;
        }
        Action::Remove => {
            remove_image()?;
        }
    }

    Ok(())
}

fn get_image_list() -> Result<Vec<Image>, Box<dyn std::error::Error>> {
    let res = reqwest::blocking::get(format!("{}/images", API_URL))?;

    if res.status().is_success() {
        let mut body: Vec<Image> = res.json()?;
        if body.is_empty() {
            return Err("No images found".into());
        }

        body.sort();
        return Ok(body);
    }

    Err(res.text()?.into())
}

fn list_images() -> Result<(), Box<dyn std::error::Error>> {
    let images = get_image_list()?;

    print_image_names(&images)?;
    return Ok(());
}

fn print_image_names(images: &Vec<Image>) -> Result<(), Box<dyn std::error::Error>> {
    for (i, image) in images.iter().enumerate() {
        let decoded_name = decode(&image.name)?;
        println!("{}. {}", i, str::from_utf8(&decoded_name)?);
    }
    Ok(())
}

fn add_image(path: &PathBuf) -> Result<(), Box<dyn std::error::Error>> {
    if !path.exists() {
        return Err("File does not exist".into());
    }

    println!("What's it called?");
    let mut title = String::new();
    io::stdin().read_line(&mut title)?;
    let encoded_title = encode(&title.trim());


    let form = Form::new()
        .text("title", encoded_title)
        .file("image", &path)?;

    let client = Client::new();

    let res = client
        .post(format!("{}/add", API_URL))
        .multipart(form)
        .send()?;

    if res.status().is_success() {
        println!("Image saved successfully");

        println!("Delete image from computer?");
        let mut answer = String::new();
        io::stdin().read_line(&mut answer)?;
        if answer.trim() == "y" || answer.trim() == "yes" {
            std::fs::remove_file(&path)?;
            println!("Image deleted");
        }

        return Ok(());
    }

    Err(res.text()?.into())
}

fn remove_image() -> Result<(), Box<dyn std::error::Error>> {
    let images = get_image_list()?;
    println!("Which image do you want to remove?");
    print_image_names(&images)?;

    let mut image_to_remove = String::new();
    io::stdin().read_line(&mut image_to_remove)?;

    let idx = image_to_remove.trim().parse::<usize>()?;
    if idx >= images.len() {
        return Err("That number doesn't have an image associated to it".into());
    }
    let image_to_remove = &images[idx].name;

    let client = Client::new();

    let res = client.delete(format!("{}/remove/{}", API_URL, image_to_remove.replace(" ", "_"))).send()?;

    if res.status().is_success() {
        println!("Image removed successfully");
        return Ok(());
    }

    Err(res.text()?.into())
}

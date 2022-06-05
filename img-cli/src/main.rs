use clap::Parser;
use serde::Deserialize;

#[derive(Parser)]
struct Args {
    #[clap(subcommand)]
    action: Action,
}

#[derive(clap::Subcommand)]
enum Action {
    List,
}

#[derive(Deserialize)]
struct Image {
    name: String,
}

fn main() {
    let args = Args::parse();
    match args.action {
        Action::List => {
            list_images();
        }
    }
}

fn list_images() {
    let body: Vec<Image> = reqwest::blocking::get("http://localhost:8000/images")
        .expect("Failed to get images")
        .json()
        .expect("Failed to parse images json");

    for image in body {
        println!("{}", image.name);
    }
}

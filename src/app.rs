use leptos::prelude::*;
use leptos::logging::log;
use leptos_meta::{provide_meta_context, MetaTags, Stylesheet, Title};
use leptos_router::{
    components::{Route, Router, Routes},
    StaticSegment,
};

use server_fn::codec::{MultipartData, MultipartFormData};
use wasm_bindgen::JsCast;
use web_sys::{FormData, HtmlFormElement, SubmitEvent};

pub fn shell(options: LeptosOptions) -> impl IntoView {
    view! {
        <!DOCTYPE html>
        <html lang="en">
            <head>
                <meta charset="utf-8"/>
                <meta name="viewport" content="width=device-width, initial-scale=1"/>
                <AutoReload options=options.clone() />
                <HydrationScripts options/>
                <MetaTags/>
            </head>
            <body>
                <App/>
            </body>
        </html>
    }
}

#[component]
pub fn App() -> impl IntoView {
    // Provides context that manages stylesheets, titles, meta tags, etc.
    provide_meta_context();

    let var_name = view! {
        // injects a stylesheet into the document <head>
        // id=leptos means cargo-leptos will hot-reload this stylesheet
        <Stylesheet id="leptos" href="/pkg/schlingel.css"/>

        // sets the document title
        <Title text="Welcome to Leptos"/>

        // content for this welcome page
        <Router>
            <main>
                <Routes fallback=|| "Page not found.".into_view()>
                    <Route path=StaticSegment("") view=crate::pages::HomePage/>
                    <Route path=StaticSegment("list") view=ListPage/>
                </Routes>
            </main>
        </Router>
    };
    var_name
}

#[component]
fn ListPage() -> impl IntoView {
    let items = (1..=5).collect::<Vec<usize>>();
    view! {
        <ul>
        {items.iter()
            .map(|r| view! { <li><ListItem index=*r></ListItem></li> })
                .collect_view()
        }
        </ul>
    }
}

#[component]
fn ListItem(index: usize) -> impl IntoView {
    let (s, _) = signal(index);
    let rsc = Resource::new(
        move || s.get(),
        // every time `count` changes, this will run
        |count| get_content(count),
    );
    view! {
        <Suspense fallback=move || view! { <span>Loading...</span> }>
        {move || Suspend::new(async move {
            let a = rsc.await;
            view! {
                <span>{a.unwrap_or(0)}</span>
            }
        })}
        </Suspense>
    }
}

#[derive(serde::Serialize, serde::Deserialize, Debug, Clone)]
struct HeftyData {
    my_file: String,
}

#[server]
pub async fn get_content(index: usize) -> Result<usize, ServerFnError> {
    tokio::time::sleep(tokio::time::Duration::from_secs(index.try_into().unwrap())).await;
    return Ok(index);
}

#[component]
pub fn FileUpload() -> impl IntoView {
    /// A simple file upload function, which does just returns the length of the file.
    ///
    /// On the server, this uses the `multer` crate, which provides a streaming API.
    #[server(
        input = MultipartFormData,
    )]
    pub async fn file_length(
        data: MultipartData,
    ) -> Result<usize, ServerFnError> {
        // `.into_inner()` returns the inner `multer` stream
        // it is `None` if we call this on the client, but always `Some(_)` on the server, so is safe to
        // unwrap
        let mut data = data.into_inner().unwrap();

        // this will just measure the total number of bytes uploaded
        let mut count = 0;
        while let Ok(Some(mut field)) = data.next_field().await {
            println!("\n[NEXT FIELD]\n");
            let name = field.name().unwrap_or_default().to_string();
            println!("  [NAME] {name}");
            while let Ok(Some(chunk)) = field.chunk().await {
                let len = chunk.len();
                count += len;
                println!("      [CHUNK] {len}");
                // in a real server function, you'd do something like saving the file here
            }
        }

        Ok(count)
    }

    let upload_action = Action::new_local(|data: &FormData| {
        // `MultipartData` implements `From<FormData>`
        file_length(data.clone().into())
    });

    view! {
        <h3>File Upload</h3>
        <p>Uploading files is fairly easy using multipart form data.</p>
        <form on:submit=move |ev: SubmitEvent| {
            ev.prevent_default();
            let target = ev.target().unwrap().unchecked_into::<HtmlFormElement>();
            let form_data = FormData::new_with_form(&target).unwrap();
            upload_action.dispatch_local(form_data);
        }>
            <input type="file" name="file_to_upload"/>
            <input type="submit"/>
        </form>
        <p>
            {move || {
                if upload_action.input_local().read().is_none() && upload_action.value().read().is_none()
                {
                    "Upload a file.".to_string()
                } else if upload_action.pending().get() {
                    "Uploading...".to_string()
                } else if let Some(Ok(value)) = upload_action.value().get() {
                    value.to_string()
                } else {
                    format!("{:?}", upload_action.value().get())
                }
            }}

        </p>
    }
}

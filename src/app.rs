use leptos::prelude::*;
use leptos_meta::{provide_meta_context, MetaTags, Stylesheet, Title};
use leptos_router::{
    components::{Route, Router, Routes},
    StaticSegment,
};

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

#[server]
pub async fn get_content(index: usize) -> Result<usize, ServerFnError> {
    crate::db::run_migrations();
    tokio::time::sleep(tokio::time::Duration::from_secs(index.try_into().unwrap())).await;
    return Ok(index);
}

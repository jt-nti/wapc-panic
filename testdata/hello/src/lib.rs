extern crate wapc_guest as guest;

use guest::prelude::*;

wapc_handler!(handle_wapc);

pub fn handle_wapc(operation: &str, msg: &[u8]) -> CallResult {
    match operation {
        "hello" => hello_world(msg),
        _ => Err("bad dispatch".into()),
    }     
}

fn hello_world(msg: &[u8]) -> CallResult {
    guest::console_log(&format!(
        "Received message: {}",
        std::str::from_utf8(msg).unwrap()
    ));
    let _res = host_call("wapc", "hello", "echo", msg)?;
    Ok(_res)
}

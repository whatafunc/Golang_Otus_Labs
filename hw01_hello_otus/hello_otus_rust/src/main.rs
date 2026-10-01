fn main() {
    let mut greeting_msg = String::from("Hello, OTUS!");
    greeting_msg = reverse::string(greeting_msg); // Reverse the order of characters in the greeting message
    println!("{}", greeting_msg); // Output the result
}

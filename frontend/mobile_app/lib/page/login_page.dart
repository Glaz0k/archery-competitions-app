import 'package:flutter/material.dart';
import 'package:keyboard_dismisser/keyboard_dismisser.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({super.key});

  @override
  State<LoginPage> createState() => _LoginPage();
}

class _LoginPage extends State<LoginPage> {
  final _sizeBlackText = const TextStyle(fontSize: 12, color: Colors.black);
  final _sizeWhiteText = const TextStyle(fontSize: 15, color: Colors.white);

  bool isPasswordVisible = true;

  final myController = TextEditingController();
  var text = '';

  @override
  void dispose() {
    myController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return KeyboardDismisser(
      gestures: [GestureType.onTap],
      child: GestureDetector(
        child: Scaffold(
          appBar: AppBar(
            title: Text('Вход'),
            backgroundColor: Colors.green,
            centerTitle: true,
          ),
          body: Center(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: <Widget>[
                SizedBox(
                  width: 300,
                  child: TextFormField(
                    controller: myController,
                    decoration: InputDecoration(
                      label: Text(
                        "Логин",
                        style: TextStyle(color: Colors.green, fontSize: 20),
                      ),
                      prefixIcon: const Icon(Icons.person),
                      hintText: 'Введите логин',
                      enabledBorder: OutlineInputBorder(
                        borderSide: BorderSide(color: Colors.teal, width: 2.2),
                        borderRadius: BorderRadius.circular(15),
                      ),
                      focusedBorder: OutlineInputBorder(
                        borderSide: BorderSide(color: Colors.teal, width: 2.2),
                        borderRadius: BorderRadius.circular(15),
                      ),
                    ),
                    validator: (String? value) {},
                    onChanged: null,
                  ),
                ),
                Container(
                  width: 300,
                  padding: EdgeInsets.only(top: 10.0),
                  child: TextFormField(
                    decoration: InputDecoration(
                      label: Text(
                        "Пароль",
                        style: TextStyle(color: Colors.green, fontSize: 20),
                      ),
                      hintText: 'Введите пароль',
                      prefixIcon: const Icon(Icons.person),
                      enabledBorder: OutlineInputBorder(
                        borderSide: BorderSide(color: Colors.teal, width: 2.2),
                        borderRadius: BorderRadius.circular(15),
                      ),
                      focusedBorder: OutlineInputBorder(
                        borderSide: BorderSide(color: Colors.teal, width: 2.2),
                        borderRadius: BorderRadius.circular(15),
                      ),
                      suffixIcon: IconButton(
                        onPressed:
                            () => setState(
                              () => isPasswordVisible = !isPasswordVisible,
                            ),
                        icon:
                            isPasswordVisible
                                ? Icon(Icons.visibility_off)
                                : Icon(Icons.visibility),
                      ),
                    ),

                    obscureText: isPasswordVisible,
                    style: _sizeBlackText,
                  ),
                ),
                Padding(
                  padding: EdgeInsets.only(top: 25.0),
                  child: ElevatedButton(
                    onPressed: () {
                      Navigator.pushNamed(context, '/profile_page');
                    },
                    style: ButtonStyle(
                      minimumSize: WidgetStateProperty.all(const Size(200, 40)),
                      backgroundColor: WidgetStateProperty.all(Colors.green),
                      overlayColor: WidgetStateProperty.all(Colors.blueAccent),
                    ),
                    child: Text("Войти", style: _sizeWhiteText),
                  ),
                ),
                ElevatedButton(
                  onPressed:
                      () => Navigator.pushNamed(context, "/competition_page"),
                  child: Text("Перейти на страницу соревнований"),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}

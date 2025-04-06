import 'package:flutter/material.dart';

class SeriesInputPage extends StatefulWidget {
  const SeriesInputPage({super.key});

  final String title = 'Результаты';

  @override
  State<SeriesInputPage> createState() => _SeriesInputPageState();
}

class _SeriesInputPageState extends State<SeriesInputPage> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        backgroundColor: Theme.of(context).colorScheme.inversePrimary,
        title: Text(widget.title),
      ),
      body: Placeholder(),
      //floatingActionButton: FloatingActionButton(
      //  onPressed: _incrementCounter,
      //  tooltip: 'Increment',
      //  child: const Icon(Icons.add),
      //),
    );
  }
}

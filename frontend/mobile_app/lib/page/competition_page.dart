import 'package:flutter/material.dart';

class CompetitionPage extends StatefulWidget {
  const CompetitionPage({super.key});

  final String title = 'Турниры и соревнования';

  @override
  State<CompetitionPage> createState() => _CompetitionPageState();
}

class _CompetitionPageState extends State<CompetitionPage> {
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

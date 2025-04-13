import 'package:flutter/material.dart';

import '../model/series_model.dart';

class SeriesPage extends StatelessWidget {
  const SeriesPage({super.key});

  final String title = 'Серии';

  @override
  Widget build(BuildContext context) {
    throw UnimplementedError();
  }
}

class SeriesCard extends StatefulWidget {
  final Widget name;
  final Series series;
  final bool active;
  final void Function() onAdd;
  final void Function(int, int) onChange;

  const SeriesCard({
    super.key,
    required this.name,
    required this.series,
    required this.onAdd,
    required this.onChange,
    required this.active,
  });

  @override
  State<StatefulWidget> createState() {
    return _SeriesCardState();
  }
}

class _SeriesCardState extends State<SeriesCard> {
  int? _selected;

  @override
  Widget build(BuildContext context) {
    List<Widget> body = [
      widget.name,
      ...widget.series.scores.indexed.map((r) {
        var (index, score) = r;
        var isSelected = index == _selected;
        var text = Text('$score');
        if (isSelected) {
          return FilledButton(onPressed: () => _selectScore(null), child: text);
        } else {
          return FilledButton.tonal(
            onPressed: () => _selectScore(index),
            child: text,
          );
        }
      }),
      FloatingActionButton(onPressed: widget.onAdd, child: Icon(Icons.add)),
    ];
    var selected = _selected;
    return Column(
      children: [
        ListView(scrollDirection: Axis.horizontal, children: body),
        AnimatedContainer(duration: Duration(microseconds: 200),
          child: selected != null? Slider(
            value: widget.series.scores[selected].toDouble(),
            divisions: 10,
            onChanged: (score) => widget.onChange(selected, score.round()),
          ) : SizedBox.shrink(),
        )
      ],
    );
  }

  void _selectScore(int? index) => setState(() {
    _selected = index;
  });
}

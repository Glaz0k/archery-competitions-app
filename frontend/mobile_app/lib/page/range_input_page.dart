import 'package:flutter/material.dart';
import 'package:mobile_app/page/widgets/onion_bar.dart';
import 'package:provider/provider.dart';

import '../api/common.dart';
import '../api/responses.dart';
import '../model/range_model.dart';

class SeriesPage extends StatelessWidget {
  const SeriesPage({super.key});

  @override
  Widget build(BuildContext context) {
    var rangeModel = context.watch<RangeModel>();
    var ranges = rangeModel.rangeGroup.ranges;
    return Scaffold(
      appBar: OnionBar("Серии", context),
      body: RefreshIndicator(
        onRefresh: () => rangeModel.reloadRangeGroup(),
        child: ListView(
          scrollDirection: Axis.vertical,
          children: [
            for (var (index, range)
                in ranges.takeWhile((range) => range.shots != null).indexed)
              RangeCard(
                name: Text("Серия ${index + 1}"),
                range: range,
                isActive: range.isActive,
              ),
          ],
        ),
      ),
    );
  }
}

class RangeCard extends StatefulWidget {
  final Widget name;
  final Range range;
  final bool isActive;

  const RangeCard({
    super.key,
    required this.name,
    required this.range,
    required this.isActive,
  });

  @override
  State<RangeCard> createState() => _RangeCardState();
}

class _RangeCardState extends State<RangeCard> {
  int? _selected;
  bool _canFinish = false;

  @override
  Widget build(BuildContext context) {
    var selected = _selected;
    var selectedShot = selected == null ? null : widget.range.shots?[selected];
    var model = Provider.of<RangeModel>(context, listen: false);
    var type = model.rangeGroup.type;
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(8.0),
        child: Column(
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                SizedBox(
                  width: 80,
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      widget.name,
                      Text("Счёт: ${widget.range.rangeScore ?? "..."}"),
                    ],
                  ),
                ),
                Expanded(
                  child: SizedBox(
                    height: 50,
                    child: ListView(
                      scrollDirection: Axis.horizontal,
                      children: [
                        for (var (index, shot) in widget.range.shots!.indexed)
                          RangeCardButton(
                            score: shot.score ?? "?",
                            isSelected: index == _selected && widget.isActive,
                            onPressed:
                                widget.isActive
                                    ? () => setState(() {
                                      if (index == _selected) {
                                        _selected = null;
                                      } else {
                                        _selected = index;
                                      }
                                    })
                                    : null,
                          ),
                      ],
                    ),
                  ),
                ),
                IconButton(
                  onPressed: () => model.registerAndEndRange(widget.range.id),
                  icon: Icon(Icons.check),
                ),
              ],
            ),
            if (selectedShot != null && widget.isActive)
              Slider(
                value: selectedShot.toSliderValue(type),
                min: type == RangeType.six2ten ? 5.0 : 0.0,
                max: 11.0,
                divisions: (type == RangeType.six2ten ? 11 - 5 : 11),
                onChanged:
                    (value) =>
                        _changeShot(widget.range, selectedShot, value, type),
              ),
          ],
        ),
      ),
    );
  }

  void _changeShot(
    Range range,
    Shot shot,
    double sliderValue,
    RangeType type,
  ) => setState(() {
    if (shot.score == null && !_canFinish) {
      _canFinish = !(range.shots?.any((s) => s.score == null) ?? true);
    }
    shot.changeBySlider(sliderValue, type);
  });
}

class RangeCardButton extends StatelessWidget {
  final String score;
  final bool isSelected;
  final VoidCallback? onPressed;

  const RangeCardButton({
    super.key,
    required this.score,
    required this.isSelected,
    required this.onPressed,
  });

  @override
  Widget build(BuildContext context) {
    var button = isSelected ? FilledButton.new : FilledButton.tonal;
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 2.0),
      child: AspectRatio(
        aspectRatio: 1,
        child: button(
          onPressed: onPressed,
          child: Text(score, textAlign: TextAlign.center),
          style: FilledButton.styleFrom(padding: EdgeInsets.all(0)),
        ),
      ),
    );
  }
}
